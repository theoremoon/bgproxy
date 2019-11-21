package main

//go:generate protoc --proto_path=../../ --go_out=plugins=grpc:../../pb bgproxy.proto
import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/theoremoon/bgproxy/common"
	"github.com/theoremoon/bgproxy/constant"
	"github.com/theoremoon/bgproxy/pb"
	"google.golang.org/grpc"
)

var (
	helpStr = fmt.Sprintf(`
pbproxyctl
  Usage:
  - pbproxyctl green -addr http://localhost:8888/ -stop "docker stop green_server"
  - pbproxyctl green -addr http://localhost:8888/ -cmd "php -S localhost:8888 -t public/"
  - pbproxyctl status
  - pbproxyctl rollback

  Commands:
  - green:    set or replace the green server
    Options:
    - addr     (required)   address of green server to check the its health
    - cmd      (optional)   command to start a green server. if use this option, stop is ignored
    - stop     (optional)   command to stop the green server when replaced or rolled back
    - status   (optional)   expected http status code. default is 200
    - limit    (optional)   maximum unhealthy limit. default is 5
    - interval (optional)   health check interval (seconds). default is 5
    - wait     (optional)   time to wait for replacing green into blue. default is 3600

  - rollback: cancel the green server and continue to use blue server
  - status:   get the current proxy status

  Global Options:
  - sock   (optional)           bgproxy server listening socket. default is %s
`, constant.Sock)
)

func help() {
	fmt.Println(helpStr)
}

func dial(sock string) (*grpc.ClientConn, error) {
	sock_split := strings.SplitN(sock, ":", 2)
	if len(sock_split) != 2 {
		return nil, errors.New("sock must follow the format: <protocol>:<address>")
	}

	dialer := func(a string, t time.Duration) (net.Conn, error) {
		return net.Dial(sock_split[0], a)
	}
	return grpc.Dial(sock_split[1], grpc.WithInsecure(), grpc.WithDialer(dialer))
}

func getStatus() error {
	set := flag.NewFlagSet("getstatusl", flag.ExitOnError)
	sock := set.String("sock", constant.Sock, "bgproxy server litening socket")
	if err := set.Parse(os.Args[2:]); err != nil {
		return err
	}

	conn, err := dial(*sock)
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewBGProxyServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	r, err := c.GetStatus(ctx, &pb.Empty{})
	if err != nil {
		return err
	}
	fmt.Println(r.GetMsg())
	return nil
}

func rollback() error {
	set := flag.NewFlagSet("rollback", flag.ExitOnError)
	sock := set.String("sock", constant.Sock, "bgproxy server litening socket")
	if err := set.Parse(os.Args[2:]); err != nil {
		return err
	}

	conn, err := dial(*sock)
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewBGProxyServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	r, err := c.Rollback(ctx, &pb.Empty{})
	if err != nil {
		return err
	}
	fmt.Println(r.GetMsg())
	return nil
}

func setGreen() error {
	// flags
	set := flag.NewFlagSet("green", flag.ExitOnError)
	addr := set.String("addr", "", "address of green server to check the its health")
	cmd := set.String("cmd", "", "command to start a green server. if use this option, stop is ignored")
	stop := set.String("stop", "", "command to stop the green server when replaced or rolled back")
	status := set.Int("status", http.StatusOK, "address of green server to check the its health")
	limit := set.Int("limit", 5, "maximum unhealthy limit")
	interval := set.Int("interval", 5, "expected http status code")
	wait := set.Int("wait", 3600, "time to wait for replacing green into blue")
	sock := set.String("sock", constant.Sock, "bgproxy server litening socket")

	if err := set.Parse(os.Args[2:]); err != nil {
		return err
	}

	conn, err := dial(*sock)
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewBGProxyServiceClient(conn)

	// If cmd is specified, run command in the background
	// and set stop command to SIGKILL
	if *cmd != "" {
		p, err := common.RunInBackground(*cmd)
		if err != nil {
			return err
		}
		*stop = fmt.Sprintf("kill -15 %d", p.Pid)
	}
	if *addr == "" || *stop == "" {
		help()
		return nil
	}

	// Send Command
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	r, err := c.SetGreen(ctx, &pb.Target{
		Url:                 *addr,
		ExpectedStatus:      int32(*status),
		UnhealthyLimit:      int32(*limit),
		WaitingTime:         int32(*wait),
		HealthcheckInterval: int32(*interval),
		StopCommand:         *stop,
	})
	if err != nil {
		return err
	}

	fmt.Println(r.GetMsg())
	return nil
}

func main() {
	if len(os.Args) == 1 {
		help()
		return
	}
	if os.Args[1] == "rollback" {
		err := rollback()
		if err != nil {
			log.Fatal(err)
		}
	} else if os.Args[1] == "status" {
		err := getStatus()
		if err != nil {
			log.Fatal(err)
		}
	} else if os.Args[1] == "green" {
		err := setGreen()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		help()
		return
	}
}
