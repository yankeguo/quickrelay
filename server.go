package quickrelay

import (
	"net"
)

type ServerOptions struct {
}

type Server interface {
	Run(l net.Listener) error
}

type server struct{}

func NewServer(opts *ServerOptions) Server {
	return &server{}
}

func (s *server) Run(l net.Listener) error {
	return nil
}
