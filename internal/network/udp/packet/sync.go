package packet

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Sync struct {
	SequentialID uint32 // 4 bytes
	Data         []byte // 1024 bytes
}

func SyncPacketSizeWithHeaders() int {
	return sequentialIDSize + dataSize
}

func ExtractSyncPacketFromBuffer(buf []byte, totalBytes int) Sync {
	return Sync{
		SequentialID: binary.BigEndian.Uint32(buf[0:4]),
		Data:         buf[4:totalBytes],
	}
}

type SyncInput struct {
	Conn net.Conn
	Pkt  Sync
}

func DispatchSync(in SyncInput) error {
	buf := make([]byte, 4+len(in.Pkt.Data))

	binary.BigEndian.PutUint32(buf[0:4], in.Pkt.SequentialID)
	copy(buf[4:], in.Pkt.Data)

	if _, err := in.Conn.Write(buf); err != nil {
		return fmt.Errorf("error dispatching ACK: %s", err.Error())
	}

	return nil
}
