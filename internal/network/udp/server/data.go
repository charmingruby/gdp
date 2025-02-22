package server

func (s *Server) receiveData() {
	// for {
	// 	pktBuffer := make([]byte, packet.DataAckPacketSizeWithHeaders())
	// 	totalBytes, _, err := s.Conn.ReadFromUDP(pktBuffer)
	// 	if err != nil {
	// 		logger.Response(fmt.Sprintf("unable to read data packet from client: %s", err.Error()))
	// 		break
	// 	}

	// 	dataAckPkt := packet.ExtractDataAckPacketFromBuffer(pktBuffer, totalBytes)

	// 	logger.Response(
	// 		fmt.Sprintf("received data packet with ack=%d, seqID=%d", dataAckPkt.AckID, dataAckPkt.SequentialID),
	// 	)

	// 	if err := packet.DispatchAck(packet.AckInput{
	// 		Conn: s.Conn,
	// 		Pkt: packet.Ack{
	// 			AckID: dataAckPkt.AckID,
	// 			Data:  dataAckPkt.Data,
	// 		},
	// 	}); err != nil {
	// 		logger.Response(fmt.Sprintf("unable to send data ack packet: %s", err.Error()))
	// 		break
	// 	}

	// 	logger.Response(
	// 		fmt.Sprintf("sent data ack packet with ack=%d, seqID=%d", dataAckPkt.AckID, dataAckPkt.SequentialID),
	// 	)
	// }
}
