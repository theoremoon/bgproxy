package main

//go:generate protoc --proto_path=../../ --go_out=plugins=grpc:../../pb bgproxy.proto
import (
	"bufio"
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

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

func (t *target) Check() bool {
	r, err := http.Get(t.Url.String())
	if err != nil {
		return false
	}
	if r.StatusCode != t.ExpectedStatus {
		return false
	}
	return true
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
	cancel context.Context
}

func (s *service) Deploy() {
	if s.Green == nil {
		return
	}
	deployTicker := time.NewTicker(s.Green.WaitToDeploy)
	defer deployTicker.Stop()
	ticker := time.NewTicker(s.Green.CheckInterval)
	defer ticker.Stop()

	unhealthyCount := 0

rollback:
	for {
		select {
		case <-deployTicker.C:
			// do deploy and exit
			logger.Println("Replacing...")
			go s.Blue.Stop()
			s.Lock()
			s.Blue = s.Green
			s.Green = nil
			s.cancel = nil
			s.Unlock()
			return

		case <-ticker.C:
			// check the health
			if s.Green.Check() {
				unhealthyCount = 0
			} else {
				logger.Println("unhealthy")
				unhealthyCount++
				if unhealthyCount >= s.Green.UnhealthyLimit {
					break rollback
				}
			}
		case <-s.cancel.Done():
			// cancel
			logger.Println("cancelled")
			break rollback
		}
	}

	// roll back
	logger.Println("Rolling back...")
	go s.Green.Stop()
	s.Lock()
	s.Green = nil
	s.cancel = nil
	s.Unlock()
}

func (s *service) SetGreen(ctx context.Context, req *pb.Target) (*pb.Result, error) {
	if s.Green != nil {
		logger.Println("to be implemented")
	}
	url, err := url.Parse(req.GetUrl())
	if err != nil {
		return nil, err
	}
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
	s.cancel = context.Background()
	s.Unlock()
	go s.Deploy()

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
	s.cancel.Done()
	return &pb.Result{
		Msg: "Rolling back...",
	}, nil
}

func run() error {
	// parse command line
	addr := flag.String("addr", "", "the host and port address to listen and serve")
	blueaddr := flag.String("blue", "", "the url to listen and serve")
	bluestop := flag.String("stop", "", "how to stop the blue server")
	grpc_socket := flag.String("sock", constant.Sock, "socket listening for gRPC server")
	flag.Parse()

	if *addr == "" {
		return errors.New("the argument [addr] is required")
	}
	if *blueaddr == "" {
		return errors.New("the argument [blue] is required")
	}
	if *bluestop == "" {
		return errors.New("the argument [stop] is required")
	}

	// split grpc_socket into grpc_net, grpc_addr
	grpc_split := strings.SplitN(*grpc_socket, ":", 2)
	if len(grpc_split) != 2 {
		return errors.New("-sock option must follow the format <network:address>")
	}

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
		service.Lock()
		if service.Green != nil {
			request.URL.Scheme = service.Green.Url.Scheme
			request.URL.Host = service.Green.Url.Host
		} else {
			request.URL.Scheme = service.Blue.Url.Scheme
			request.URL.Host = service.Blue.Url.Host
		}
		service.Unlock()
	}
	rp := httputil.ReverseProxy{
		Director: director,
	}
	server := http.Server{
		Addr:    *addr,
		Handler: &rp,
	}

	// grpc server
	g_server := grpc.NewServer()
	pb.RegisterBGProxyServiceServer(g_server, service)
	reflection.Register(g_server)

	conn, err := net.Listen(grpc_split[0], grpc_split[1])
	if err != nil {
		return err
	}

	// serve
	err_ch := make(chan error)
	go func() {
		err_ch <- server.ListenAndServe()
	}()
	go func() {
		err_ch <- g_server.Serve(conn)
	}()
	return <-err_ch
}

func main() {
	logger.Fatal(run())
}
