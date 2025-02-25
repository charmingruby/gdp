package server

import (
	"fmt"
	"math/rand/v2"

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

		s.adjustWindowSize()

		ackPkt := packet.Ack{
			AckID:      dataAckPkt.SequentialID + 1,
			Data:       dataAckPkt.Data,
			WindowSize: s.windowSize,
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
			fmt.Sprintf("sent ack packet with ack=%d, windowSize=%d", ackPkt.AckID, ackPkt.WindowSize),
		)

		currentServerSequentialID++
		currentClientSequentialID++
	}
}

func (s *Server) adjustWindowSize() {
	randomLoss := 0.01 + rand.Float64()*(1.00-0.01)

	if float64(s.congestionControl.PackageLoss) > randomLoss {
		logger.HighlightedErrorResponse(fmt.Sprintf("handling error: Packet Loss Detected, currentWindowSize=%d, cwnd=%d, sshthresh=%d", s.windowSize, s.congestionControl.Cwnd, s.congestionControl.Sshthresh))
		s.congestionControl.Sshthresh = s.congestionControl.Cwnd / 2
		s.congestionControl.Cwnd = 1
	} else if s.congestionControl.Cwnd < s.congestionControl.Sshthresh {
		logger.HighlightedErrorResponse(fmt.Sprintf("handling error: Exponential Window Growth, currentWindowSize=%d, cwnd=%d, sshthresh=%d", s.windowSize, s.congestionControl.Cwnd, s.congestionControl.Sshthresh))
		s.congestionControl.Cwnd *= 2
	} else {
		logger.HighlightedErrorResponse(fmt.Sprintf("handling error: Linear Window Growth, currentWindowSize=%d, cwnd=%d, sshthresh=%d", s.windowSize, s.congestionControl.Cwnd, s.congestionControl.Sshthresh))
		s.congestionControl.Cwnd++
	}

	logger.HighlightedErrorResponse(fmt.Sprintf("error handled, new values: currentWindowSize=%d, cwnd=%d, sshthresh=%d", s.windowSize, s.congestionControl.Cwnd, s.congestionControl.Sshthresh))

	s.windowSize = s.congestionControl.Cwnd
}
