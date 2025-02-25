package client

import (
	"fmt"

	"github.com/charmingruby/gdp/internal/network/udp/client/messaging/extract"
	"github.com/charmingruby/gdp/internal/network/udp/client/messaging/packet"
	"github.com/charmingruby/gdp/pkg/logger"
)

func (c *Client) sendData(serverSequentialID, clientSequentialID uint32) {
	var currentServerSequentialID uint32 = serverSequentialID + 1
	var currentClientSequentialID uint32 = clientSequentialID + 1

	for range c.packageLoadSize {
		ackPktBuf := make([]byte, packet.AckPacketSizeWithHeaders())
		ackPkt := packet.DataAck{
			AckID:        currentServerSequentialID,
			SequentialID: currentClientSequentialID,
			Data:         ackPktBuf,
		}

		if err := packet.DispatchDataAck(packet.DataAckInput{
			Conn: c.Conn,
			Pkt:  ackPkt,
		}); err != nil {
			logger.Response(fmt.Sprintf("unable to send last synchronization packet: %s", err.Error()))
			continue
		}

		logger.Response(
			fmt.Sprintf("sent data ack packet with ack=%d, seq=%d", ackPkt.AckID, ackPkt.SequentialID),
		)

		ackBuf := extract.NewAckSyncBuffer()
		_, err := c.Conn.Read(ackBuf)
		if err != nil {
			logger.Response(fmt.Sprintf("error receiving ack packet: %s", err.Error()))
			continue
		}

		receivedAckPkt := extract.AckPacketFromBuffer(ackBuf)

		c.allowedWindow = receivedAckPkt.WindowSize

		logger.Response(
			fmt.Sprintf("ack packet received: ack=%d, windowSize=%d", receivedAckPkt.AckID, receivedAckPkt.WindowSize),
		)

		currentClientSequentialID++
		currentServerSequentialID++
		c.allowedWindow--
	}
}
