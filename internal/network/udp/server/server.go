package server

import (
	"fmt"
	"net"

	"github.com/charmingruby/gdp/internal/network/udp/shared/congestion"
	"github.com/charmingruby/gdp/pkg/logger"
)

type CongestionThreshold struct {
	PackageLoss float32
}

type ServerInput struct {
	Port      int
	Threshold CongestionThreshold
}

type Server struct {
	Conn              *net.UDPConn
	addr              *net.UDPAddr
	windowSize        uint32
	congestionControl congestion.Congestion
}

func New(in ServerInput) (*Server, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", in.Port))
	if err != nil {
		return nil, fmt.Errorf("unable to resolve UDP address: %s", err.Error())
	}

	return &Server{
		addr:       addr,
		windowSize: 10,
		congestionControl: congestion.Congestion{
			PackageLoss: in.Threshold.PackageLoss,
			Cwnd:        1,
			Sshthresh:   16,
		},
	}, nil
}

func (s *Server) Listen() error {
	conn, err := net.ListenUDP("udp", s.addr)
	if err != nil {
		return fmt.Errorf("unable to listen on UDP address: %s", err.Error())
	}

	s.Conn = conn

	return nil
}

func (s *Server) Read() error {
	var baseSequentialID = uint32(0)

	logger.Header("Synchronization Process")
	logger.OpenBracket()

	syncResult, err := s.sync(baseSequentialID)
	if err != nil {
		return err
	}

	logger.Response(fmt.Sprintf(
		"synchronization completed: serverSequentialID=%d, clientSequentialID=%d",
		syncResult.serverSequentialID,
		syncResult.clientSequentialID,
	))

	logger.CloseBracket()

	logger.Divider()

	logger.Header("Data Transfer Process")
	logger.OpenBracket()

	s.receiveData(syncResult.serverSequentialID, syncResult.clientSequentialID)

	logger.CloseBracket()

	return nil
}
