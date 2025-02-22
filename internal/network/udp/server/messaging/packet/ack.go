package packet

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/charmingruby/gdp/internal/network/udp/shared/constant"
)

type Ack struct {
	AckID uint32 // 4 bytes
	Data  []byte // 1024 bytes
}

func AckPacketSizeWithHeaders() int {
	return constant.ACK_ID_SIZE + constant.DATA_SIZE
}

type AckInput struct {
	Conn       *net.UDPConn
	ClientAddr *net.UDPAddr
	Pkt        Ack
}

func DispatchAck(in AckInput) error {
	ackBuffer := make([]byte, 8+len(in.Pkt.Data))
	binary.BigEndian.PutUint32(ackBuffer[0:4], in.Pkt.AckID)
	copy(ackBuffer[4:], in.Pkt.Data)

	if _, err := in.Conn.WriteToUDP(ackBuffer, in.ClientAddr); err != nil {
		return fmt.Errorf("error dispatching ACK: %s", err.Error())
	}

	return nil
}
