package server

import (
	"fmt"

	"github.com/charmingruby/gdp/internal/network/udp/server/messaging/extract"
	"github.com/charmingruby/gdp/internal/network/udp/server/messaging/packet"
	"github.com/charmingruby/gdp/pkg/logger"
)

func (s *Server) termination(currentSequentialID uint32) error {
	logger.Response("waiting for fin packet from client...")

	finPktBuf := extract.NewFinBuffer()
	_, clientAddr, err := s.Conn.ReadFromUDP(finPktBuf)
	if err != nil {
		return fmt.Errorf("unable to read fin packet from client: %s", err.Error())
	}

	finPkt := extract.FinPacketFromBuffer(finPktBuf)

	logger.Response(
		fmt.Sprintf("received fin packet with fin=%d", finPkt.Fin),
	)

	ackFinPkt := packet.AckFin{
		AckID: finPkt.Fin + 1,
		Fin:   uint32(currentSequentialID),
	}

	if err := packet.DispatchAckFin(packet.AckFinInput{
		Conn:       s.Conn,
		ClientAddr: clientAddr,
		Pkt:        ackFinPkt,
	}); err != nil {
		return fmt.Errorf("unable to send ack-fin packet: %s", err.Error())
	}

	logger.Response(
		fmt.Sprintf("sent ack-fin packet with ack=%d, fin=%d", ackFinPkt.AckID, ackFinPkt.Fin),
	)

	ackPktBuffer := extract.NewAckBuffer()
	totalBytes, _, err := s.Conn.ReadFromUDP(ackPktBuffer)
	if err != nil {
		return fmt.Errorf("unable to read ack packet from client: %s", err.Error())
	}

	ackPkt := extract.AckPacketFromBuffer(ackPktBuffer, totalBytes)

	logger.Response(
		fmt.Sprintf("received ack packet with ack=%d", ackPkt.AckID),
	)

	return nil
}
