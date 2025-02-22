package server

import (
	"fmt"

	"github.com/charmingruby/gdp/internal/network/udp/server/messaging/extract"
	"github.com/charmingruby/gdp/internal/network/udp/server/messaging/packet"
	"github.com/charmingruby/gdp/pkg/logger"
)

func (s *Server) receiveData(serverSequentialID, clientSequentialID uint32) {
	var currentServerSequentialID uint32 = serverSequentialID + 1
	var currentClientSequentialID uint32 = clientSequentialID + 1

	for {
		dataPktBuf := extract.NewDataAckBuffer()
		totalBytes, clientAddr, err := s.Conn.ReadFromUDP(dataPktBuf)
		if err != nil {
			logger.Response(fmt.Sprintf("unable to read data packet from client: %s", err.Error()))
			continue
		}

		dataAckPkt := extract.DataAckPacketFromBuffer(dataPktBuf, totalBytes)

		logger.Response(
			fmt.Sprintf("received data packet with ack=%d, seqID=%d", dataAckPkt.AckID, dataAckPkt.SequentialID),
		)

		ackPkt := packet.Ack{
			AckID: dataAckPkt.SequentialID + 1,
			Data:  dataAckPkt.Data,
		}

		if err := packet.DispatchAck(packet.AckInput{
			Conn:       s.Conn,
			ClientAddr: clientAddr,
			Pkt:        ackPkt,
		}); err != nil {
			logger.Response(fmt.Sprintf("unable to send data ack packet: %s", err.Error()))
			continue
		}

		logger.Response(
			fmt.Sprintf("sent ack packet with ack=%d", ackPkt.AckID),
		)

		currentServerSequentialID++
		currentClientSequentialID++
	}
}
