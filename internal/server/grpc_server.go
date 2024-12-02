package server

import (
	"net"

	"github.com/ex0rcist/metflix/internal/entities"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	address entities.Address
	server  *grpc.Server
	notify  chan error
}

func NewGRPCServer(server *grpc.Server, address entities.Address) *GRPCServer {
	s := &GRPCServer{
		address: address,
		server:  server,
		notify:  make(chan error, 1),
	}

	return s
}

func (s *GRPCServer) Start() {
	go func() {
		listen, err := net.Listen("tcp", s.address.String())
		if err != nil {
			s.notify <- err
			return
		}

		s.notify <- s.server.Serve(listen)
		close(s.notify)
	}()
}

func (s *GRPCServer) Notify() <-chan error {
	return s.notify
}

func (s *GRPCServer) Shutdown() {
	if s.server == nil {
		return
	}

	s.server.GracefulStop()
}
