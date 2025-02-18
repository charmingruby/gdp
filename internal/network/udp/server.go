package udp

import (
	"fmt"
	"net"

	"github.com/charmingruby/gdp/internal/network/udp/congestion"
	"github.com/charmingruby/gdp/internal/network/udp/packet"
	"github.com/charmingruby/gdp/internal/shared/logger"
)

type CongestionThreshold struct {
	PackageLoss float32
}

type ServerInput struct {
	Port      int
	Threshold CongestionThreshold
}

type Server struct {
	Conn       *net.UDPConn
	addr       *net.UDPAddr
	congestion congestion.Congestion
}

func NewServer(in ServerInput) (*Server, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", in.Port))
	if err != nil {
		return nil, fmt.Errorf("unable to resolve UDP address: %s", err.Error())
	}

	return &Server{
		addr:       addr,
		congestion: congestion.Congestion(in.Threshold),
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
	var serverSequentialID = uint32(0)

	logger.Header("Synchronization Process")
	logger.OpenBracket()

	if err := s.sync(serverSequentialID); err != nil {
		return err
	}
	logger.CloseBracket()

	logger.Divider()

	return nil
}

func (s *Server) sync(baseSequentialID uint32) error {

	logger.Response("waiting for sync packet from client...")

	pktBuffer := make([]byte, packet.AckSyncPacketSizeWithHeaders())
	totalBytes, clientAddr, err := s.Conn.ReadFromUDP(pktBuffer)
	if err != nil {
		return fmt.Errorf("unable to read synchronize packet from client: %s", err.Error())
	}

	syncPkt := packet.ExtractSyncPacketFromBuffer(pktBuffer, totalBytes)

	logger.Response(
		fmt.Sprintf("received synchronize packet with sequentialID=%d", syncPkt.SequentialID),
	)

	ackSyncPkt := packet.AckSync{
		AckID:        syncPkt.SequentialID + 1,
		SequentialID: baseSequentialID,
		Data:         make([]byte, packet.DataSize()),
	}

	if err := packet.DispatchAckSync(packet.AckSyncInput{
		Conn:       s.Conn,
		ClientAddr: clientAddr,
		Pkt:        ackSyncPkt,
	}); err != nil {
		return fmt.Errorf("unable to send ack-sync packet: %s", err.Error())
	}

	logger.Response(
		fmt.Sprintf("sent ack-sync packet with ack=%d, seqID=%d", ackSyncPkt.AckID, ackSyncPkt.SequentialID),
	)

	pktBuffer = make([]byte, packet.AckSyncPacketSizeWithHeaders())
	totalBytes, _, err = s.Conn.ReadFromUDP(pktBuffer)
	if err != nil {
		return fmt.Errorf("unable to read synchronize packet from client: %s", err.Error())
	}

	ackPkt := packet.ExtractAckPacketFromBuffer(pktBuffer, totalBytes)

	logger.Response(
		fmt.Sprintf("received last synchronization ack packet with ack=%d", ackPkt.AckID),
	)

	return nil
}
