package extract

import (
	"encoding/binary"

	"github.com/charmingruby/gdp/internal/network/udp/client/messaging/packet"
)

func NewSyncBuffer() []byte {
	return make([]byte, packet.SyncPacketSizeWithHeaders())
}

func SyncPacketFromBuffer(buf []byte, totalBytes int) packet.Sync {
	return packet.Sync{
		SequentialID: binary.BigEndian.Uint32(buf[0:4]),
		Data:         buf[4:totalBytes],
	}
}
