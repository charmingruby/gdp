package client

import (
	"fmt"

	"github.com/charmingruby/gdp/internal/network/udp/client/messaging/packet"
	"github.com/charmingruby/gdp/pkg/logger"
)

func (c *Client) sendData(serverSequentialID, clientSequentialID uint32) error {
	var currentServerSequentialID uint32 = serverSequentialID
	var currentClientSequentialID uint32 = clientSequentialID

	for range c.config.PackageLoadSize {
		ackPktBuf := make([]byte, packet.AckPacketSizeWithHeaders())
		ackPkt := packet.DataAck{
			AckID:        currentServerSequentialID + 1,
			SequentialID: currentClientSequentialID + 1,
			Data:         ackPktBuf,
		}

		if err := packet.DispatchDataAck(packet.DataAckInput{
			Conn: c.Conn,
			Pkt:  ackPkt,
		}); err != nil {
			return fmt.Errorf("unable to send last synchronization packet: %s", err.Error())
		}

		logger.Response(
			fmt.Sprintf("ack packet with ack=%d, seq=%d", ackPkt.AckID, ackPkt.SequentialID),
		)
	}

	return nil
}
