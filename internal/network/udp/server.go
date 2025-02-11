package udp

import (
	"fmt"
	"net"
)

type ServerInput struct {
	Port      int
	Threshold Threshold
}

type Server struct {
	Conn *net.UDPConn

	threshold Threshold
	addr      *net.UDPAddr
}

func NewServer(in ServerInput) (*Server, error) {
	addr, err := net.ResolveUDPAddr(
		"udp",
		fmt.Sprintf(":%d", in.Port),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve UDP address: %s", err.Error())
	}

	return &Server{
		addr:      addr,
		threshold: in.Threshold,
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

func (s *Server) Read() {
	pktBuffer := make([]byte, defaultPacketSize())
	expectedSequentialID := uint32(0)

	for {
		totalBytes, clientAddr, err := s.Conn.ReadFromUDP(pktBuffer)
		if err != nil {
			continue
		}

		pkt := extractPacketFromBuffer(pktBuffer, totalBytes)

		if isOcurrence := isAPackageLossOccurence(s.threshold.PackageLoss); isOcurrence {
			fmt.Printf("Package loss ocurred for package with sequential ID %d\n", pkt.SequentialID)
			continue
		}

		isPackageOrdered := pkt.SequentialID == expectedSequentialID
		if isPackageOrdered {
			fmt.Printf("Received package with sequential ID %d\n", pkt.SequentialID)
			expectedSequentialID++
		} else {
			fmt.Printf("Received UNORDERED package with sequential ID %d\n", pkt.SequentialID)
		}

		if err := dispatchAck(ackInput{
			conn:                 s.Conn,
			clientAddr:           clientAddr,
			pkt:                  pkt,
			expectedSequentialID: expectedSequentialID,
		}); err != nil {
			fmt.Println(err.Error())
		}
	}
}
