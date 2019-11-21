package main

//go:generate protoc --proto_path=../../ --go_out=plugins=grpc:../../pb bgproxy.proto
import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/theoremoon/bgproxy/common"
	"github.com/theoremoon/bgproxy/constant"
	"github.com/theoremoon/bgproxy/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	logger = log.New(os.Stderr, "bgproxy:", log.Flags())
)

type target struct {
	Url            *url.URL
	ExpectedStatus int
	UnhealthyLimit int
	CheckInterval  time.Duration
	WaitToDeploy   time.Duration
	StopCommand    string
	DeployedAt     time.Time
}

func (t *target) Check() (bool, string) {
	r, err := http.Get(t.Url.String())
	if err != nil {
		return false, err.Error()
	}
	if r.StatusCode != t.ExpectedStatus {
		return false, fmt.Sprintf("Returned status code: %d", r.StatusCode)
	}
	return true, ""
}

/// Stop stops the target
func (t *target) Stop() error {
	if t.StopCommand == "" {
		return nil
	}

	cmd := exec.Command("sh", "-c", t.StopCommand)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		s := bufio.NewScanner(stdout)
		for s.Scan() {
			logger.Println(s.Text())
		}
	}()
	go func() {
		s := bufio.NewScanner(stderr)
		for s.Scan() {
			logger.Println(s.Text())
		}
	}()
	return cmd.Wait()
}

type service struct {
	sync.Mutex
	Blue   *target
	Green  *target
	cancel *context.CancelFunc
}

func (s *service) Deploy(ctx context.Context) {
	if s.Green == nil {
		return
	}
	deployTicker := time.NewTicker(s.Green.WaitToDeploy)
	defer deployTicker.Stop()
	ticker := time.NewTicker(s.Green.CheckInterval)
	defer ticker.Stop()

	unhealthyCount := 0
	logger.Println("Green registered")

rollback:
	for {
		select {
		case <-deployTicker.C:
			// do deploy and exit
			logger.Println("Replacing...")
			logger.Println("Stopping blue...")
			go s.Blue.Stop()
			s.Lock()
			s.Blue = s.Green
			s.Green = nil
			s.cancel = nil
			s.Unlock()
			logger.Println("Done")
			return

		case <-ticker.C:
			// check the health
			if ok, reason := s.Green.Check(); ok {
				unhealthyCount = 0
			} else {
				logger.Println("Unhealthy: " + reason)
				unhealthyCount++
				if unhealthyCount >= s.Green.UnhealthyLimit {
					break rollback
				}
			}
		case <-ctx.Done():
			// cancel
			logger.Println("Cancelled")
			break rollback
		}
	}

	// roll back
	logger.Println("Rolling back...")
	logger.Println("Stopping green...")
	go s.Green.Stop()
	s.Lock()
	s.Green = nil
	s.cancel = nil
	s.Unlock()
	logger.Println("Done")
}

func (s *service) SetGreen(ctx context.Context, req *pb.Target) (*pb.Result, error) {
	if s.Green != nil {
		return &bp.Result{
			Msg: "Green exists. First roll it back",
		}, nil
	}
	url, err := url.Parse(req.GetUrl())
	if err != nil {
		return nil, err
	}
	deployCtx, cancel := context.WithCancel(context.Background())

	s.Lock()
	s.Green = &target{
		Url:            url,
		ExpectedStatus: int(req.GetExpectedStatus()),
		UnhealthyLimit: int(req.GetUnhealthyLimit()),
		CheckInterval:  time.Duration(req.GetHealthcheckInterval()) * time.Second,
		WaitToDeploy:   time.Duration(req.GetWaitingTime()) * time.Second,
		StopCommand:    req.GetStopCommand(),
		DeployedAt:     time.Now(),
	}
	s.cancel = &cancel
	s.Unlock()
	go s.Deploy(deployCtx)

	return &pb.Result{
		Msg: "OK",
	}, nil
}
func (s *service) Rollback(ctx context.Context, req *pb.Empty) (*pb.Result, error) {
	if s.cancel == nil {
		return &pb.Result{
			Msg: "Green doesn't running",
		}, nil
	}
	(*s.cancel)()
	return &pb.Result{
		Msg: "Rolling back...",
	}, nil
}

var (
	helpStr = fmt.Sprintf(`
bgproxy
  Usage:
  - bgproxy -addr localhost:9999 -blue http://localhost:8888/ -stop "docker stop blue"

  Options:
  - addr  (required)  address to be listened by this proxy
  - blue  (required)  initial blue url to pass the request
  - cmd   (selective) command to start blue server. this command runs in background.
  - stop  (selective) how to stop the blue server when replaced it with green.
  - sock  (optional)  socket listening by gRPC server. default is %s
`, constant.Sock)
)

func help() {
	fmt.Println(helpStr)
}

func run() error {
	// parse command line
	addr := flag.String("addr", "", "the host and port address to listen and serve")
	bluecmd := flag.String("cmd", "", "command to start blue server")
	blueaddr := flag.String("blue", "", "the url to listen and serve")
	bluestop := flag.String("stop", "", "how to stop the blue server")
	grpc_socket := flag.String("sock", constant.Sock, "socket listening for gRPC server")
	flag.Parse()

	if *addr == "" || *blueaddr == "" || (*bluecmd == "" && *bluestop == "") {
		help()
		return nil
	}

	// split grpc_socket into grpc_net, grpc_addr
	grpc_split := strings.SplitN(*grpc_socket, ":", 2)
	if len(grpc_split) != 2 {
		return errors.New("-sock option must follow the format <network:address>")
	}

	// if cmd is specified, launch it and set stop command
	p, err := common.RunInBackground(*bluecmd)
	if err != nil {
		return err
	}
	*bluestop = fmt.Sprintf("kill -9 %d", p.Pid)

	// blue-green
	url, err := url.Parse(*blueaddr)
	if err != nil {
		return err
	}
	service := &service{
		Blue: &target{
			Url:         url,
			StopCommand: *bluestop,
			DeployedAt:  time.Now(),
		},
		Green:  nil,
		cancel: nil,
	}

	// reverse proxy
	director := func(request *http.Request) {
		target := service.Green
		if target == nil {
			target = service.Blue
		}

		service.Lock()
		if target.Url.Scheme == "unix" {
			request.URL.Scheme = "http" // dummy
			request.URL.Host = "socket" // dummy
		} else {
			request.URL.Scheme = target.Url.Scheme
			request.URL.Host = target.Url.Host
		}
		service.Unlock()
	}
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			service.Lock()
			defer service.Unlock()

			target := service.Green
			if target == nil {
				target = service.Blue
			}
			if target.Url.Scheme == "unix" {
				return net.Dial("unix", target.Url.Path)
			} else {
				return net.Dial("tcp", target.Url.Host)
			}
		},
	}
	rp := httputil.ReverseProxy{
		Director:  director,
		Transport: transport,
	}

	// serve (graceful shutdown (to close unix domain socket))
	err_ch := make(chan error)
	sig_ch := make(chan os.Signal)
	signal.Notify(sig_ch, os.Interrupt)

	var listener net.Listener
	if strings.HasPrefix(*addr, "unix:") {
		listener, err = net.Listen("unix", strings.TrimPrefix(*addr, "unix:"))
		if err != nil {
			return err
		}
		defer listener.Close()
	}

	go func() {
		server := http.Server{
			Handler: &rp,
		}

		if listener != nil {
			err_ch <- server.Serve(listener)
		} else {
			server.Addr = *addr
			err_ch <- server.ListenAndServe()
		}
	}()

	conn, err := net.Listen(grpc_split[0], grpc_split[1])
	if err != nil {
		return err
	}
	defer conn.Close()
	go func() {
		// grpc server
		g_server := grpc.NewServer()
		pb.RegisterBGProxyServiceServer(g_server, service)
		reflection.Register(g_server)

		err_ch <- g_server.Serve(conn)
	}()
	logger.Println("Start")

	select {
	case err := <-err_ch:
		return err
	case <-sig_ch:
		if service.cancel != nil {
			(*service.cancel)()
		}
		return nil
	}
}

func main() {
	err := run()
	if err != nil {
		logger.Fatal(err)
	}
	// when run returns nil, it may caused by SIGINT
	os.Exit(130)
}
