package udp

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net"
	"time"
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
	var base uint32 = 0
	var nextSequentialID uint32 = 0
	var cwnd int = c.clientThreshold.InitialWindowSize
	var ssthresh int = c.clientThreshold.InitialSshthresh
	var dupAckCount int = 0
	var lastAck uint32 = 0

	for i := 0; i < 2; i++ {
		var window []packet
		for nextSequentialID < base+uint32(cwnd) {
			data := make([]byte, packetSize())
			rand.Read(data)
			packet := packet{SequentialID: nextSequentialID, Data: data}
			window = append(window, packet)
			nextSequentialID++
		}

		for _, pkt := range window {
			buf := make([]byte, 8+len(pkt.Data))
			binary.BigEndian.PutUint32(buf[0:4], pkt.SequentialID)
			binary.BigEndian.PutUint32(buf[4:8], pkt.AckID)
			copy(buf[8:], pkt.Data)
			_, err := c.Conn.Write(buf)
			if err != nil {
				return fmt.Errorf("unable to send packet: %s", err.Error())
			}

			fmt.Printf("Packet sent: SequentialID=%d, cwnd=%d, ssthresh=%d\n", pkt.SequentialID, cwnd, ssthresh)
		}

		timeout := time.Duration(c.clientThreshold.TimeoutInSeconds) * time.Second

		ackChan := make(chan uint32)
		go func() {
			buf := make([]byte, 8)
			c.Conn.SetReadDeadline(time.Now().Add(timeout))
			_, err := c.Conn.Read(buf)
			if err != nil {
				fmt.Println("Error receiving ACK:", err)
				return
			}
			ackNum := binary.BigEndian.Uint32(buf[4:8])
			ackChan <- ackNum
		}()

		select {
		case ackNum := <-ackChan:
			fmt.Printf("ACK received: %d\n", ackNum)

			if ackNum == lastAck {
				dupAckCount++
				if dupAckCount == 3 {
					fmt.Println("Fast Retransmit: 3 duplicated ACKs  detected")
					ssthresh = cwnd / 2
					cwnd = ssthresh + 3
					base = ackNum + 1
					nextSequentialID = base
					dupAckCount = 0
					fmt.Printf("Fast Retransmit: cwnd=%d, ssthresh=%d\n", cwnd, ssthresh)
				}
			} else {
				dupAckCount = 0
				if cwnd < ssthresh {
					// Slow Start
					cwnd *= 2
				} else {
					// Congestion Avoidance
					cwnd++
				}
				base = ackNum + 1
			}
			lastAck = ackNum

		case <-time.After(timeout):
			fmt.Println("Timeout detected")
			ssthresh = cwnd / 2
			cwnd = int(c.clientThreshold.InitialWindowSize)
			base = lastAck + 1
			nextSequentialID = base
			fmt.Printf("Timeout: cwnd=%d, ssthresh=%d\n", cwnd, ssthresh)
		}
	}

	return nil
}
