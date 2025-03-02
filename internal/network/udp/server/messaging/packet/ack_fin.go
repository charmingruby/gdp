package packet

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/charmingruby/gdp/internal/network/udp/shared/constant"
)

type AckFin struct {
	AckID uint32 // 4 bytes
	Fin   uint32 // 4 bytes
}

func AckFinPacketSizeWithHeaders() int {
	return constant.ACK_ID_SIZE + constant.SEQUENTIAL_ID_SIZE
}

type AckFinInput struct {
	Conn       *net.UDPConn
	ClientAddr *net.UDPAddr
	Pkt        AckFin
}

func DispatchAckFin(in AckFinInput) error {
	finBuffer := make([]byte, 8)

	binary.BigEndian.PutUint32(finBuffer[0:4], in.Pkt.AckID)
	binary.BigEndian.PutUint32(finBuffer[4:8], in.Pkt.Fin)

	if _, err := in.Conn.WriteToUDP(finBuffer, in.ClientAddr); err != nil {
		return fmt.Errorf("error dispatching Fin: %s", err.Error())
	}

	return nil
}
