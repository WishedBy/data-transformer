package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Server struct {
	Port      int
	TLSConfig *tls.Config
	server    *http.Server
}

func (server *Server) Run() (err error) {
	finish := make(chan error)
	if server.Port == 0 {
		server.Port = 443
	}

	server.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", server.Port),
		Handler:           server,
		ReadHeaderTimeout: 10 * time.Second,
	}

	//Start the http server
	go server.startServer(finish)

	return <-finish
}

//Starts a HTTPS server
func (server *Server) startServer(finish chan error) {

	var ln net.Listener
	var err error

	ln, err = net.Listen("tcp", server.server.Addr)
	if err != nil {
		<-finish
		return
	}

	if server.TLSConfig != nil {
		ln = tls.NewListener(ln, server.TLSConfig)
	}

	finish <- server.server.Serve(ln)
}

func (server *Server) ServeHTTP(response http.ResponseWriter, request *http.Request) {

}
