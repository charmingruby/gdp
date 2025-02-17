package packet

import (
	"encoding/binary"
	"fmt"
	"net"
)

type AckSync struct {
	AckID        uint32 // 4 bytes
	SequentialID uint32 // 4 bytes
	Data         []byte // 1024 bytes
}

func AckSyncPacketSizeWithHeaders() int {
	return ackIDPacketSize + sequentialIDSize + dataSize
}

func ExtractAckSyncPacketFromBuffer(buf []byte, totalBytes int) AckSync {
	return AckSync{
		AckID:        binary.BigEndian.Uint32(buf[0:4]),
		SequentialID: binary.BigEndian.Uint32(buf[4:8]),
		Data:         buf[8:totalBytes],
	}
}

type AckSyncInput struct {
	Conn       *net.UDPConn
	ClientAddr *net.UDPAddr
	Pkt        AckSync
}

func DispatchAckSync(in AckSyncInput) error {
	ackBuffer := make([]byte, 8+len(in.Pkt.Data))
	binary.BigEndian.PutUint32(ackBuffer[0:4], in.Pkt.AckID)
	binary.BigEndian.PutUint32(ackBuffer[4:8], in.Pkt.SequentialID)
	copy(ackBuffer[8:], in.Pkt.Data)

	if _, err := in.Conn.WriteToUDP(ackBuffer, in.ClientAddr); err != nil {
		return fmt.Errorf("error dispatching ACK+sync: %s", err.Error())
	}

	return nil
}
