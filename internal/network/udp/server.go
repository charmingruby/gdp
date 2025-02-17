package udp

import (
	"fmt"
	"net"

	"github.com/charmingruby/gdp/internal/network/udp/packet"
)

type ServerInput struct {
	Port      int
	Threshold CongestionThreshold
}

type Server struct {
	Conn *net.UDPConn

	threshold CongestionThreshold
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

func (s *Server) Read() error {
	var serverSequentialID = uint32(0)
	var clientSequentialID = uint32(0)
	var clientAddr *net.UDPAddr = nil
	var isSync = false

	for {
		if !isSync {
			fmt.Printf("waiting for sync packet from client...\n")

			pktBuffer := make([]byte, packet.AckPacketSizeWithHeaders())

			totalBytes, incomeClientAddr, err := s.Conn.ReadFromUDP(pktBuffer)
			if err != nil {
				return fmt.Errorf("unable to read syncronize packet from client: %s", err.Error())
			}

			syncPkt := packet.ExtractSyncPacketFromBuffer(pktBuffer, totalBytes)

			fmt.Printf("received syncronize packet with sequentialID=%d\n", syncPkt.SequentialID)

			clientSequentialID = syncPkt.SequentialID
			clientAddr = incomeClientAddr
			isSync = true
		} else {
			pktBuffer := make([]byte, packet.AckPacketSizeWithHeaders())

			totalBytes, incomeClientAddr, err := s.Conn.ReadFromUDP(pktBuffer)
			if err != nil {
				continue
			}

			if clientAddr == nil {
				clientAddr = incomeClientAddr
			}

			pkt := packet.ExtractSyncPacketFromBuffer(pktBuffer, totalBytes)

			if isOcurrence := isAPackageLossOccurence(s.threshold.PackageLoss); isOcurrence {
				fmt.Printf("package loss ocurred for package with sequential ID %d\n", pkt.SequentialID)
				continue
			}

			clientSequentialID = pkt.SequentialID + 1
		}

		pkt := packet.AckSync{
			AckID:        clientSequentialID + 1,
			SequentialID: serverSequentialID,
			Data:         make([]byte, packet.DataSize()),
		}

		if err := packet.DispatchAckSync(packet.AckSyncInput{
			Conn:       s.Conn,
			ClientAddr: clientAddr,
			Pkt:        pkt,
		}); err != nil {
			fmt.Println(err.Error())
		}
	}
}
