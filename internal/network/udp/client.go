package udp

import (
	"fmt"
	"net"

	"github.com/charmingruby/gdp/internal/network/udp/packet"
	"github.com/charmingruby/gdp/internal/shared/logger"
)

type ClientInput struct {
	ServerPort      int
	ClientThreshold ClientThreshold
}

type ClientThreshold struct {
	TimeoutInSeconds  int
	InitialWindowSize int
	MaxWindowSize     int
	InitialSshthresh  int
}

type Client struct {
	Conn net.Conn

	serverAddr      string
	clientThreshold ClientThreshold
}

func NewClient(in ClientInput) Client {
	return Client{
		serverAddr:      fmt.Sprintf(":%d", in.ServerPort),
		clientThreshold: in.ClientThreshold,
	}
}

func (c *Client) Run() error {
	conn, err := net.Dial("udp", c.serverAddr)
	if err != nil {
		return fmt.Errorf("unable to connect on UDP address: %s", err.Error())
	}

	c.Conn = conn

	return nil
}

func (c *Client) Dispatch() error {
	var baseSequentialID uint32 = 10

	logger.Header("Synchronization Process")
	logger.OpenBracket()

	_, err := c.sync(baseSequentialID)
	if err != nil {
		return err
	}

	logger.CloseBracket()
	logger.Divider()

	// for {

	// 	buf := make([]byte, 8+len(pkt.Data))
	// 	binary.BigEndian.PutUint32(buf[0:4], pkt.SequentialID)
	// 	binary.BigEndian.PutUint32(buf[4:8], pkt.AckID)
	// 	copy(buf[8:], pkt.Data)

	// 	_, err := c.Conn.Write(buf)
	// 	if err != nil {
	// 		return fmt.Errorf("unable to send packet: %s", err.Error())
	// 	}

	// 	fmt.Printf("Packet sent: SequentialID=%d, cwnd=%d, ssthresh=%d", pkt.SequentialID, cwnd, ssthresh)

	// 	timeout := time.Duration(c.clientThreshold.TimeoutInSeconds) * time.Second

	// 	ackChan := make(chan uint32)
	// 	go func() {
	// 		buf := make([]byte, 8)
	// 		c.Conn.SetReadDeadline(time.Now().Add(timeout))
	// 		_, err := c.Conn.Read(buf)
	// 		if err != nil {
	// 			fmt.Println("Error receiving ACK:", err)
	// 			return
	// 		}
	// 		ackNum := binary.BigEndian.Uint32(buf[4:8])
	// 		ackChan <- ackNum
	// 	}()

	// 	select {
	// 	case ackNum := <-ackChan:
	// 		fmt.Printf("ACK received: %d", ackNum)

	// 		if ackNum == lastAck {
	// 			dupAckCount++
	// 			if dupAckCount == 3 {
	// 				fmt.Println("Fast Retransmit: 3 duplicated ACKs  detected")
	// 				ssthresh = cwnd / 2
	// 				cwnd = ssthresh + 3
	// 				base = ackNum + 1
	// 				nextSequentialID = base
	// 				dupAckCount = 0
	// 				fmt.Printf("Fast Retransmit: cwnd=%d, ssthresh=%d", cwnd, ssthresh)
	// 			}
	// 		} else {
	// 			dupAckCount = 0
	// 			if cwnd < ssthresh {
	// 				// Slow Start
	// 				cwnd *= 2
	// 			} else {
	// 				// Congestion Avoidance
	// 				cwnd++
	// 			}
	// 			base = ackNum + 1
	// 		}
	// 		lastAck = ackNum

	// 	case <-time.After(timeout):
	// 		fmt.Println("Timeout detected")
	// 		ssthresh = cwnd / 2
	// 		cwnd = int(c.clientThreshold.InitialWindowSize)
	// 		base = lastAck + 1
	// 		nextSequentialID = base
	// 		fmt.Printf("Timeout: cwnd=%d, ssthresh=%d", cwnd, ssthresh)
	// 	}
	// }

	return nil
}

func (c *Client) sync(baseSequentialID uint32) (packet.Ack, error) {
	syncPktBuf := make([]byte, packet.SyncPacketSizeWithHeaders())
	syncPkt := packet.Sync{
		SequentialID: baseSequentialID,
		Data:         syncPktBuf,
	}

	if err := packet.DispatchSync(packet.SyncInput{
		Conn: c.Conn,
		Pkt:  syncPkt,
	}); err != nil {
		return packet.Ack{}, fmt.Errorf("unable to send sync packet: %s", err.Error())
	}

	logger.Response(
		fmt.Sprintf("sent sync packet with sequentialID=%d", syncPkt.SequentialID),
	)

	ackSyncReceiverBuf := make([]byte, packet.AckSyncPacketSizeWithHeaders())
	_, err := c.Conn.Read(ackSyncReceiverBuf)
	if err != nil {
		return packet.Ack{}, fmt.Errorf("error receiving sync-ack packet: %s", err.Error())
	}

	ackSyncPkt := packet.ExtractAckSyncPacketFromBuffer(ackSyncReceiverBuf, packet.AckSyncPacketSizeWithHeaders())

	logger.Response(
		fmt.Sprintf("ack-sync packet received: ack=%d", ackSyncPkt.AckID),
	)

	ackPktBuf := make([]byte, packet.AckPacketSizeWithHeaders())
	ackPkt := packet.Ack{
		AckID: ackSyncPkt.SequentialID + 1,
		Data:  ackPktBuf,
	}

	if err := packet.DispatchAck(packet.AckInput{
		Conn: c.Conn,
		Pkt:  ackPkt,
	}); err != nil {
		return packet.Ack{}, fmt.Errorf("unable to send last synchronization packet: %s", err.Error())
	}

	logger.Response(
		fmt.Sprintf("sent last ack packet with ack=%d", ackPkt.AckID),
	)

	return ackPkt, nil
}
