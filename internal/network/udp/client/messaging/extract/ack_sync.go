package extract

import (
	"encoding/binary"

	"github.com/charmingruby/gdp/internal/network/udp/server/messaging/packet"
)

func NewAckSyncBuffer() []byte {
	return make([]byte, packet.AckSyncPacketSizeWithHeaders())
}

func AckSyncPacketFromBuffer(buf []byte) packet.AckSync {
	return packet.AckSync{
		AckID:        binary.BigEndian.Uint32(buf[0:4]),
		SequentialID: binary.BigEndian.Uint32(buf[4:8]),
		Data:         buf[8:packet.AckSyncPacketSizeWithHeaders()],
	}
}
