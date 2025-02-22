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
		AckID: binary.BigEndian.Uint32(buf[0:4]),
		Data:  buf[4:packet.AckPacketSizeWithHeaders()],
	}
}
