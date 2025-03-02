package extract

import (
	"encoding/binary"

	"github.com/charmingruby/gdp/internal/network/udp/server/messaging/packet"
)

func NewAckFinBuffer() []byte {
	return make([]byte, packet.AckFinPacketSizeWithHeaders())
}

func AckFinPacketFromBuffer(buf []byte) packet.AckFin {
	return packet.AckFin{
		AckID: binary.BigEndian.Uint32(buf[0:4]),
		Fin:   binary.BigEndian.Uint32(buf[4:8]),
	}
}
