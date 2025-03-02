package packet

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/charmingruby/gdp/internal/network/udp/shared/constant"
)

type Fin struct {
	Fin uint32 // 4 bytes
}

func FinPacketSizeWithHeaders() int {
	return constant.SEQUENTIAL_ID_SIZE
}

type FinInput struct {
	Conn net.Conn
	Pkt  Fin
}

func DispatchFin(in FinInput) error {
	finBuffer := make([]byte, 8)

	binary.BigEndian.PutUint32(finBuffer[0:4], in.Pkt.Fin)

	if _, err := in.Conn.Write(finBuffer); err != nil {
		return fmt.Errorf("error dispatching Fin: %s", err.Error())
	}

	return nil
}
