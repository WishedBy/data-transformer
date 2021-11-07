package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Method string

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	PATCH  = "PATCH"
	DELETE = "DELETE"
)

type Router struct {
	http.ServeMux
	handlers map[string]map[Method]http.HandlerFunc
}

func (r *Router) createHandler(pattern string) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		handler, ok := r.handlers[pattern][Method(req.Method)]
		if !ok {
			rw.WriteHeader(404)
			return
		}
		handler(rw, req)
	}
}
func (r *Router) AddHandler(pattern string, method Method, handler http.HandlerFunc) {
	if r.handlers == nil {
		r.handlers = map[string]map[Method]http.HandlerFunc{}
	}
	if _, ok := r.handlers[pattern]; !ok {
		r.HandleFunc(pattern, r.createHandler(pattern))
		r.handlers[pattern] = map[Method]http.HandlerFunc{}
	}
	r.handlers[pattern][method] = handler
}

type Server struct {
	Port        int
	TLSConfig   *tls.Config
	server      *http.Server
	router      *Router
	handlerBuff map[string]map[Method]http.HandlerFunc
}

func (s *Server) AddHandler(pattern string, method Method, handler http.HandlerFunc) {
	if s.router != nil {
		s.router.AddHandler(pattern, method, handler)
		return
	}

	if s.handlerBuff == nil {
		s.handlerBuff = map[string]map[Method]http.HandlerFunc{}
	}
	if _, ok := s.handlerBuff[pattern]; !ok {
		s.handlerBuff[pattern] = map[Method]http.HandlerFunc{}
	}
	s.handlerBuff[pattern][method] = handler
}

func (s *Server) Run() (err error) {
	finish := make(chan error)
	if s.Port == 0 {
		s.Port = 443
	}
	s.router = &Router{ServeMux: *http.NewServeMux()}
	for p, subBuff := range s.handlerBuff {
		for m, h := range subBuff {
			s.router.AddHandler(p, m, h)
		}
	}
	s.handlerBuff = nil

	s.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.Port),
		Handler:           s.router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go s.startServer(finish)

	return <-finish
}

func (s *Server) StopServer(ctx context.Context) error {
	if s.server != nil {

	}
	return s.server.Shutdown(ctx)
}
func (s *Server) startServer(finish chan error) {

	var ln net.Listener
	var err error

	ln, err = net.Listen("tcp", s.server.Addr)
	if err != nil {
		<-finish
		return
	}

	if s.TLSConfig != nil {
		ln = tls.NewListener(ln, s.TLSConfig)
	}

	finish <- s.server.Serve(ln)
}
