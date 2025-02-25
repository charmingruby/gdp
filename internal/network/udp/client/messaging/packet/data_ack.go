package packet

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/charmingruby/gdp/internal/network/udp/shared/constant"
)

type DataAck struct {
	AckID        uint32 // 4 bytes
	SequentialID uint32 // 4 bytes
	WindowSize   uint32 // 4 bytes
	Data         []byte // 1024 bytes
}

func DataAckPacketSizeWithHeaders() int {
	return constant.ACK_ID_SIZE + constant.SEQUENTIAL_ID_SIZE + constant.DATA_SIZE
}

type DataAckInput struct {
	Conn net.Conn
	Pkt  DataAck
}

func DispatchDataAck(in DataAckInput) error {
	ackBuffer := make([]byte, 12+len(in.Pkt.Data))

	binary.BigEndian.PutUint32(ackBuffer[0:4], in.Pkt.AckID)
	binary.BigEndian.PutUint32(ackBuffer[4:8], in.Pkt.SequentialID)
	binary.BigEndian.PutUint32(ackBuffer[8:12], in.Pkt.WindowSize)
	copy(ackBuffer[12:], in.Pkt.Data)

	if _, err := in.Conn.Write(ackBuffer); err != nil {
		return fmt.Errorf("error dispatching data ACK: %s", err.Error())
	}

	return nil
}
