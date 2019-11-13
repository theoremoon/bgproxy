package main

//go:generate protoc --go_out=plugins=grpc:pb bgproxy.proto
import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/theoremoon/bgproxy/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type target struct {
	Url            *url.URL
	ExpectedStatus int
	CheckInterval  time.Duration
	WaitToDeploy   time.Duration
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

rollback:
	for {
		select {
		case <-deployTicker.C:
			// do deploy and exit
			log.Println("Replacing...")
			s.Lock()
			s.Blue = s.Green
			s.Green = nil
			s.cancel = nil
			s.Unlock()
			return

		case <-ticker.C:
			// check the health
			if s.Green.Check() == false {
				log.Println("unhealthy")
				break rollback
			}

		case <-s.cancel.Done():
			// cancel
			log.Println("cancelled")
			break rollback
		}
	}

	// roll back
	log.Println("Rolling back...")
	s.Lock()
	s.Green = nil
	s.cancel = nil
	s.Unlock()
}

func (s *service) SetGreen(ctx context.Context, req *pb.Target) (*pb.Result, error) {
	if s.Green != nil {
		log.Println("to be implemented")
	}
	url, err := url.Parse(req.GetUrl())
	if err != nil {
		return nil, err
	}
	s.Lock()
	s.Green = &target{
		Url:            url,
		ExpectedStatus: int(req.GetExpectedStatus()),
		CheckInterval:  time.Duration(req.GetHealthcheckInterval()) * time.Second,
		WaitToDeploy:   time.Duration(req.GetWaitingTime()) * time.Second,
	}
	s.cancel = context.Background()
	s.Unlock()
	go s.Deploy()

	return &pb.Result{
		Msg: "OK",
	}, nil
}

func run() error {
	// parse command line
	addr := flag.String("addr", "", "the host and port address to listen and serve")
	blueaddr := flag.String("blue", "", "the url to listen and serve")
	grpc_socket := flag.String("sock", "unix:/tmp/pbproxy.sock", "socket listening for gRPC server")
	flag.Parse()

	if *addr == "" {
		return errors.New("the argument [addr] is required")
	}
	if *blueaddr == "" {
		return errors.New("the argument [blue] is required")
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
			Url:            url,
			ExpectedStatus: http.StatusOK,
			CheckInterval:  5 * time.Second,
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
	log.Fatal(run())
}
