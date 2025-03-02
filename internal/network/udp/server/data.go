package server

import (
	"fmt"

	"github.com/charmingruby/gdp/internal/network/udp/server/messaging/extract"
	"github.com/charmingruby/gdp/internal/network/udp/server/messaging/packet"
	"github.com/charmingruby/gdp/internal/network/udp/shared/analysis"
	"github.com/charmingruby/gdp/internal/network/udp/shared/congestion"
	"github.com/charmingruby/gdp/pkg/logger"
)

func (s *Server) receiveData(serverSequentialID, clientSequentialID uint32) uint32 {
	var currentServerSequentialID uint32 = serverSequentialID
	var currentClientSequentialID uint32 = clientSequentialID
	var packagesReceived int = 0
	var packagesLost int = 0
	var retransmissionData []analysis.RetransmissionUnit

	for range s.packageLoadSize {
		dataPktBuf := extract.NewDataAckBuffer()
		totalBytes, clientAddr, err := s.Conn.ReadFromUDP(dataPktBuf)
		packagesReceived++
		if err != nil {
			logger.Response(fmt.Sprintf("unable to read data packet from client: %s", err.Error()))
			continue
		}

		dataAckPkt := extract.DataAckPacketFromBuffer(dataPktBuf, totalBytes)

		logger.Response(
			fmt.Sprintf("received data packet with ack=%d, seqID=%d", dataAckPkt.AckID, dataAckPkt.SequentialID),
		)

		packageLost := congestion.IsAPackageLossOccurence(s.congestionControl.PackageLoss)
		if packageLost {
			packagesLost++
		}

		retransmissionResult := s.adjustWindowSize(packageLost)
		retransmissionData = append(retransmissionData, retransmissionResult)

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

	if err := analysis.SaveRetransmissionData("./data/retransmission.json", retransmissionData); err != nil {
		logger.HighlightedErrorResponse(fmt.Sprintf("unable to save retransmission data: %s", err.Error()))
	}

	if err := analysis.SavePackageLossData("./data/package-loss.json", analysis.PackageLossData{
		PackagesReceived: packagesReceived,
		PackagesLost:     packagesLost,
	}); err != nil {
		logger.HighlightedErrorResponse(fmt.Sprintf("unable to save package loss data: %s", err.Error()))
	}

	return currentServerSequentialID
}

func (s *Server) adjustWindowSize(packageLost bool) analysis.RetransmissionUnit {
	if packageLost {
		logger.HighlightedErrorResponse(fmt.Sprintf("handling error: Packet Loss Detected, currentWindowSize=%d, cwnd=%d, sshthresh=%d", s.windowSize, s.congestionControl.Cwnd, s.congestionControl.Sshthresh))
		currentCwnd := s.congestionControl.Cwnd
		currentWindowSize := s.windowSize
		currentSshthresh := s.congestionControl.Sshthresh

		s.congestionControl.Sshthresh = s.congestionControl.Cwnd / 2
		s.congestionControl.Cwnd = 1
		s.windowSize = s.congestionControl.Cwnd

		return analysis.RetransmissionUnit{
			Type:             "Package Loss",
			InitialRWND:      int(currentWindowSize),
			FinalRWND:        int(s.windowSize),
			InitialCWND:      int(currentCwnd),
			FinalCWND:        int(s.congestionControl.Cwnd),
			InitialSshthresh: int(currentSshthresh),
			FinalSshthresh:   int(s.congestionControl.Sshthresh),
		}
	} else if s.congestionControl.Cwnd < s.congestionControl.Sshthresh {
		logger.HighlightedErrorResponse(fmt.Sprintf("handling error: Exponential Window Growth, currentWindowSize=%d, cwnd=%d, sshthresh=%d", s.windowSize, s.congestionControl.Cwnd, s.congestionControl.Sshthresh))
		currentCwnd := s.congestionControl.Cwnd
		currentWindowSize := s.windowSize

		s.congestionControl.Cwnd *= 2
		s.windowSize = s.congestionControl.Cwnd

		return analysis.RetransmissionUnit{
			Type:             "Exponential Window Growth",
			InitialRWND:      int(currentWindowSize),
			FinalRWND:        int(s.windowSize),
			InitialCWND:      int(currentCwnd),
			FinalCWND:        int(s.congestionControl.Cwnd),
			InitialSshthresh: int(s.congestionControl.Sshthresh),
			FinalSshthresh:   int(s.congestionControl.Sshthresh),
		}
	} else {
		logger.HighlightedErrorResponse(fmt.Sprintf("handling error: Linear Window Growth, currentWindowSize=%d, cwnd=%d, sshthresh=%d", s.windowSize, s.congestionControl.Cwnd, s.congestionControl.Sshthresh))
		currentCwnd := s.congestionControl.Cwnd
		currentWindowSize := s.windowSize

		s.congestionControl.Cwnd++
		s.windowSize = s.congestionControl.Cwnd

		return analysis.RetransmissionUnit{
			Type:             "Linear Window Growth",
			InitialRWND:      int(currentWindowSize),
			FinalRWND:        int(s.windowSize),
			InitialCWND:      int(currentCwnd),
			FinalCWND:        int(s.congestionControl.Cwnd),
			InitialSshthresh: int(s.congestionControl.Sshthresh),
			FinalSshthresh:   int(s.congestionControl.Sshthresh),
		}
	}
}
