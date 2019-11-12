package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type target struct {
	url            url.URL
	expectedStatus int
	checkInterval  time.Duration
}

func (t *target) healthCheck(unhealthy chan error) {
	ticker := time.NewTicker(t.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// check health
			r, err := http.Get(t.url.toString())
			if err != nil {
				unhealthy <- err
				return
			}
			if r.StatusCode != t.expectedStatus {
				unhealthy <- fmt.Errorf("expected status is %d but %d is returned", t.expectedStatus, r.StatusCode)
				return
			}
		}
	}

}

func run() error {
	addr := flag.String("addr", "", "the host and port address to listen and serve")
	flag.Parse()

	if *addr == "" {
		return errors.New("the argument [addr] is required")
	}

	// reverse proxy
	director := func(request *http.Request) {
		request.URL.Scheme = "http" // temporary
		request.URL.Host = ":8888"  // temporary
	}
	rp := httputil.ReverseProxy{
		Director: director,
	}

	// serve
	server := http.Server{
		Addr:    *addr,
		Handler: &rp,
	}
	return server.ListenAndServe()
}

func main() {
	log.Fatal(run())
}
