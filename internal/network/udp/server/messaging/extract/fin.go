package extract

import (
	"encoding/binary"

	"github.com/charmingruby/gdp/internal/network/udp/client/messaging/packet"
)

func NewFinBuffer() []byte {
	return make([]byte, packet.FinPacketSizeWithHeaders())
}

func FinPacketFromBuffer(buf []byte) packet.Fin {
	return packet.Fin{
		Fin: binary.BigEndian.Uint32(buf[0:4]),
	}
}
