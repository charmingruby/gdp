package udp

import "encoding/binary"

type packet struct {
	AckID        uint32 // 4 bytes
	SequentialID uint32 // 4 bytes
	Data         []byte // 1024 bytes
}

func extractPackageFromBuffer(buf []byte, totalBytes int) packet {
	return packet{
		AckID:        binary.BigEndian.Uint32(buf[0:4]),
		SequentialID: binary.BigEndian.Uint32(buf[4:8]),
		Data:         buf[8:totalBytes],
	}
}
