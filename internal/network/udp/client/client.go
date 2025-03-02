package client

import (
	"fmt"
	"net"

	"github.com/charmingruby/gdp/pkg/logger"
)

type Client struct {
	Conn net.Conn

	serverAddr      string
	packageLoadSize int
	allowedWindow   uint32
}

type ClientInput struct {
	ServerPort      int
	PackageLoadSize int
}

func New(in ClientInput) Client {
	return Client{
		serverAddr:      fmt.Sprintf(":%d", in.ServerPort),
		packageLoadSize: in.PackageLoadSize,
		allowedWindow:   0,
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

	syncResult, err := c.sync(baseSequentialID)
	if err != nil {
		return err
	}

	logger.Response(fmt.Sprintf(
		"synchronization completed: serverSequentialID=%d, clientSequentialID=%d",
		syncResult.serverSequentialID,
		syncResult.clientSequentialID,
	))

	logger.CloseBracket()
	logger.Divider()

	logger.Header("Data Transfer Process")
	logger.OpenBracket()

	currentClientSequentialID := c.sendData(syncResult.serverSequentialID, syncResult.clientSequentialID)

	logger.CloseBracket()
	logger.Divider()

	logger.Header("Termination Process")
	logger.OpenBracket()

	if err := c.termination(currentClientSequentialID); err != nil {
		logger.HighlightedErrorResponse(err.Error())
	}

	logger.CloseBracket()
	return nil
}
