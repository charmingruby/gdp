package client

import (
	"fmt"

	"github.com/charmingruby/gdp/internal/network/udp/client/messaging/extract"
	"github.com/charmingruby/gdp/internal/network/udp/client/messaging/packet"
	"github.com/charmingruby/gdp/pkg/logger"
)

type syncResult struct {
	clientSequentialID uint32
	serverSequentialID uint32
}

func (c *Client) sync(baseSequentialID uint32) (syncResult, error) {
	syncPktBuf := make([]byte, packet.SyncPacketSizeWithHeaders())
	syncPkt := packet.Sync{
		SequentialID: baseSequentialID,
		Data:         syncPktBuf,
	}

	if err := packet.DispatchSync(packet.SyncInput{
		Conn: c.Conn,
		Pkt:  syncPkt,
	}); err != nil {
		return syncResult{}, fmt.Errorf("unable to send sync packet: %s", err.Error())
	}

	logger.Response(
		fmt.Sprintf("sent sync packet with sequentialID=%d", syncPkt.SequentialID),
	)

	ackSyncReceiverBuf := extract.NewAckSyncBuffer()
	_, err := c.Conn.Read(ackSyncReceiverBuf)
	if err != nil {
		return syncResult{}, fmt.Errorf("error receiving sync-ack packet: %s", err.Error())
	}

	ackSyncPkt := extract.AckSyncPacketFromBuffer(ackSyncReceiverBuf)

	logger.Response(
		fmt.Sprintf("ack-sync packet received: ack=%d", ackSyncPkt.AckID),
	)

	ackPktBuf := make([]byte, packet.AckPacketSizeWithHeaders())
	ackPkt := packet.Ack{
		AckID: ackSyncPkt.SequentialID + 1,
		Data:  ackPktBuf,
	}

	if err := packet.DispatchAck(packet.AckInput{
		Conn: c.Conn,
		Pkt:  ackPkt,
	}); err != nil {
		return syncResult{}, fmt.Errorf("unable to send last synchronization packet: %s", err.Error())
	}

	logger.Response(
		fmt.Sprintf("sent last ack packet with ack=%d", ackPkt.AckID),
	)

	return syncResult{
		serverSequentialID: ackPkt.AckID,
		clientSequentialID: ackSyncPkt.AckID,
	}, nil
}
