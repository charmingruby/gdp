package client

import (
	"fmt"

	"github.com/charmingruby/gdp/internal/network/udp/client/messaging/extract"
	"github.com/charmingruby/gdp/internal/network/udp/client/messaging/packet"
	"github.com/charmingruby/gdp/pkg/logger"
)

func (s *Client) termination(currentClientSequentialID uint32) error {
	logger.Response("waiting for fin packet from client...")

	finPkt := packet.Fin{
		Fin: uint32(currentClientSequentialID),
	}

	if err := packet.DispatchFin(packet.FinInput{
		Conn: s.Conn,
		Pkt:  finPkt,
	}); err != nil {
		return fmt.Errorf("unable to dispatch fin packet from client: %s", err.Error())
	}

	logger.Response(
		fmt.Sprintf("sent fin packet with fin=%d", finPkt.Fin),
	)

	ackFinPktBuf := extract.NewAckFinBuffer()
	_, err := s.Conn.Read(ackFinPktBuf)
	if err != nil {
		return fmt.Errorf("unable to read ack-fin packet from client: %s", err.Error())
	}

	ackFinPkt := extract.AckFinPacketFromBuffer(ackFinPktBuf)

	logger.Response(
		fmt.Sprintf("received fin packet with fin=%d, ack=%d", ackFinPkt.Fin, ackFinPkt.AckID),
	)

	ackPkt := packet.Ack{
		AckID: ackFinPkt.Fin + 1,
		Data:  []byte{},
	}

	if err := packet.DispatchAck(packet.AckInput{
		Conn: s.Conn,
		Pkt:  ackPkt,
	}); err != nil {
		return fmt.Errorf("unable to send ack packet: %s", err.Error())
	}

	logger.Response(
		fmt.Sprintf("sent ack packet with ack=%d", ackPkt.AckID),
	)

	return nil
}
