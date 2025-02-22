package client

import (
	"fmt"
	"net"

	"github.com/charmingruby/gdp/pkg/logger"
)

type ClientInput struct {
	ServerPort int
	Config     ClientConfig
}

type ClientConfig struct {
	TimeoutInSeconds  int
	InitialWindowSize int
	MaxWindowSize     int
	InitialSshthresh  int
	PackageLoadSize   int
}

type Client struct {
	Conn net.Conn

	serverAddr string
	config     ClientConfig
}

func New(in ClientInput) Client {
	return Client{
		serverAddr: fmt.Sprintf(":%d", in.ServerPort),
		config:     in.Config,
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

	c.sendData(syncResult.serverSequentialID, syncResult.clientSequentialID)

	logger.CloseBracket()

	return nil
}
