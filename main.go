package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WishedBy/data-transformer/pkg/server"
)

var terminateCB = []func() error{}

type config struct {
	port int
}

func configure() config {
	cfg := config{}
	flag.IntVar(&cfg.port, "http-port", 80, "the port the api will run on")

	flag.Parse()
	return cfg
}

func main() {

	shutdownSigs := make(chan os.Signal, 1)
	signal.Notify(shutdownSigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		<-shutdownSigs
		for _, cb := range terminateCB {

			if err := cb(); err != nil {
				fmt.Printf("%+v\n\n", err)
			}
		}
		done <- true
	}()

	cfg := configure()
	runServer(cfg)

	<-done
}

func runServer(cfg config) *server.Server {
	srv := &server.Server{Port: cfg.port}

	srv.AddHandler("/", server.GET, func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("hello"))
	})

	terminateCB = append(terminateCB, func() error {
		return stopServer(srv)
	})

	if err := srv.Run(); err != nil {
		fmt.Printf("%+v\n\n", err)
	}
	return srv
}

func stopServer(server *server.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return server.StopServer(ctx)
}
