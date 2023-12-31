package am

import (
	"net/http"
)

type (
	Server struct {
		Core
		opts []Option
		http.Server
		router *Router
	}
)

func NewServer(name string, opts ...Option) (server *Server) {
	return &Server{
		Core:   NewCore(name, opts...),
		opts:   opts,
		router: NewRouter("root-router", "/", opts...),
	}
}
