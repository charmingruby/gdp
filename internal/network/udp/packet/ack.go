package packet

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Ack struct {
	AckID uint32 // 4 bytes
	Data  []byte // 1024 bytes
}

func AckPacketSizeWithHeaders() int {
	return ackIDPacketSize + dataSize
}

func ExtractAckPacketFromBuffer(buf []byte, totalBytes int) Ack {
	return Ack{
		AckID: binary.BigEndian.Uint32(buf[0:4]),
		Data:  buf[4:totalBytes],
	}
}

type AckInput struct {
	Conn net.Conn
	Pkt  Ack
}

func DispatchAck(in AckInput) error {
	ackBuffer := make([]byte, 4+len(in.Pkt.Data))
	binary.BigEndian.PutUint32(ackBuffer[0:4], in.Pkt.AckID)
	copy(ackBuffer[4:], in.Pkt.Data)

	if _, err := in.Conn.Write(ackBuffer); err != nil {
		return fmt.Errorf("error dispatching ACK: %s", err.Error())
	}

	return nil
}
