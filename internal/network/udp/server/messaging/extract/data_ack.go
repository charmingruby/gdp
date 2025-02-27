package extract

import (
	"encoding/binary"

	"github.com/charmingruby/gdp/internal/network/udp/client/messaging/packet"
)

func NewDataAckBuffer() []byte {
	return make([]byte, packet.DataAckPacketSizeWithHeaders())
}

func DataAckPacketFromBuffer(buf []byte, totalBytes int) packet.DataAck {
	return packet.DataAck{
		AckID:        binary.BigEndian.Uint32(buf[0:4]),
		SequentialID: binary.BigEndian.Uint32(buf[4:8]),
		WindowSize:   binary.BigEndian.Uint32(buf[8:12]),
		Data:         buf[12:totalBytes],
	}
}
