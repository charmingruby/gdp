package extract

import (
	"encoding/binary"

	"github.com/charmingruby/gdp/internal/network/udp/client/messaging/packet"
)

func DataAckPacketFromBuffer(buf []byte, totalBytes int) packet.DataAck {
	return packet.DataAck{
		AckID:        binary.BigEndian.Uint32(buf[0:4]),
		SequentialID: binary.BigEndian.Uint32(buf[4:8]),
		Data:         buf[8:totalBytes],
	}
}
