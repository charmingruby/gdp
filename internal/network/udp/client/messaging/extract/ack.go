package extract

import (
	"encoding/binary"

	"github.com/charmingruby/gdp/internal/network/udp/server/messaging/packet"
)

func NewAckBuffer() []byte {
	return make([]byte, packet.AckPacketSizeWithHeaders())
}

func AckPacketFromBuffer(buf []byte) packet.Ack {
	return packet.Ack{
		AckID:      binary.BigEndian.Uint32(buf[0:4]),
		WindowSize: binary.BigEndian.Uint32(buf[4:8]),
		Data:       buf[12:packet.AckPacketSizeWithHeaders()],
	}
}
