package udp

import "encoding/binary"

const (
	ackIDPacketSize  = 4
	sequentialIDSize = 4
	dataSize         = 1024
)

type packet struct {
	AckID        uint32 // 4 bytes
	SequentialID uint32 // 4 bytes
	Data         []byte // 1024 bytes
}

func defaultPacketSize() int {
	return ackIDPacketSize + sequentialIDSize + dataSize
}

func extractPacketFromBuffer(buf []byte, totalBytes int) packet {
	return packet{
		AckID:        binary.BigEndian.Uint32(buf[0:4]),
		SequentialID: binary.BigEndian.Uint32(buf[4:8]),
		Data:         buf[8:totalBytes],
	}
}
