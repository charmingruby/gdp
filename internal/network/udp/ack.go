package udp

import (
	"encoding/binary"
	"fmt"
	"net"
)

type ackInput struct {
	conn                 *net.UDPConn
	clientAddr           *net.UDPAddr
	pkg                  packet
	expectedSequentialID uint32
}

func dispatchAck(in ackInput) error {
	ackBuffer := make([]byte, 8)
	binary.BigEndian.PutUint32(ackBuffer[0:4], in.pkg.SequentialID)
	binary.BigEndian.PutUint32(ackBuffer[4:8], in.expectedSequentialID)

	if _, err := in.conn.WriteToUDP(ackBuffer, in.clientAddr); err != nil {
		return fmt.Errorf("error dispatching ACK: %s", err.Error())
	}

	return nil
}
