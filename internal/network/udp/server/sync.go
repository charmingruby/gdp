package server

import (
	"fmt"

	"github.com/charmingruby/gdp/internal/network/udp/server/messaging/extract"
	"github.com/charmingruby/gdp/internal/network/udp/server/messaging/packet"
	"github.com/charmingruby/gdp/internal/network/udp/shared/constant"
	"github.com/charmingruby/gdp/pkg/logger"
)

type syncResult struct {
	clientSequentialID uint32
	serverSequentialID uint32
}

func (s *Server) sync(baseSequentialID uint32) (syncResult, error) {
	logger.Response("waiting for sync packet from client...")

	syncPktBuf := extract.NewSyncBuffer()
	totalBytes, clientAddr, err := s.Conn.ReadFromUDP(syncPktBuf)
	if err != nil {
		return syncResult{}, fmt.Errorf("unable to read synchronize packet from client: %s", err.Error())
	}
	syncPkt := extract.SyncPacketFromBuffer(syncPktBuf, totalBytes)

	logger.Response(
		fmt.Sprintf("received synchronize packet with sequentialID=%d", syncPkt.SequentialID),
	)

	ackSyncPkt := packet.AckSync{
		AckID:        syncPkt.SequentialID + 1,
		SequentialID: baseSequentialID,
		Data:         make([]byte, constant.DATA_SIZE),
	}

	if err := packet.DispatchAckSync(packet.AckSyncInput{
		Conn:       s.Conn,
		ClientAddr: clientAddr,
		Pkt:        ackSyncPkt,
	}); err != nil {
		return syncResult{}, fmt.Errorf("unable to send ack-sync packet: %s", err.Error())
	}

	logger.Response(
		fmt.Sprintf("sent ack-sync packet with ack=%d, seqID=%d", ackSyncPkt.AckID, ackSyncPkt.SequentialID),
	)

	ackPktBuffer := extract.NewAckBuffer()
	totalBytes, _, err = s.Conn.ReadFromUDP(ackPktBuffer)
	if err != nil {
		return syncResult{}, fmt.Errorf("unable to read synchronize packet from client: %s", err.Error())
	}

	ackPkt := extract.AckPacketFromBuffer(ackPktBuffer, totalBytes)

	logger.Response(
		fmt.Sprintf("received last synchronization ack packet with ack=%d", ackPkt.AckID),
	)

	return syncResult{
		serverSequentialID: ackPkt.AckID,
		clientSequentialID: ackSyncPkt.AckID,
	}, nil
}
