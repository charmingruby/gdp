package udp

import (
	"fmt"
	"net"
)

type ServerInput struct {
	Port int
}

type Server struct {
	addr *net.UDPAddr
}

func NewServer(in ServerInput) (*Server, error) {
	addr, err := net.ResolveUDPAddr(
		"udp",
		fmt.Sprintf(":%d", in.Port),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve UDP address: %s", err.Error())
	}

	return &Server{addr: addr}, nil
}

func (s *Server) Start() (*net.UDPConn, error) {
	conn, err := net.ListenUDP("udp", s.addr)
	if err != nil {
		return nil, fmt.Errorf("unable to listen on UDP address: %s", err.Error())
	}

	return conn, nil
}
