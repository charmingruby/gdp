package packet

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/charmingruby/gdp/internal/network/udp/shared/constant"
)

type Sync struct {
	SequentialID uint32 // 4 bytes
	Data         []byte // 1024 bytes
}

func SyncPacketSizeWithHeaders() int {
	return constant.SEQUENTIAL_ID_SIZE + constant.DATA_SIZE
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
