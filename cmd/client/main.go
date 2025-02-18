package main

import (
	"fmt"

	"github.com/charmingruby/gdp/config"
	"github.com/charmingruby/gdp/internal/network/udp"
	"github.com/charmingruby/gdp/internal/shared/logger"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	clientCfg := cfg.Client

	client := udp.NewClient(udp.ClientInput{
		ServerPort: clientCfg.ServerPort,
		ClientThreshold: udp.ClientThreshold{
			TimeoutInSeconds:  clientCfg.ClientThreshold.TimeoutInSeconds,
			InitialWindowSize: clientCfg.ClientThreshold.InitialWindowSize,
			MaxWindowSize:     clientCfg.ClientThreshold.MaxWindowSize,
			InitialSshthresh:  clientCfg.ClientThreshold.InitialSshthresh,
		},
	})

	if err := client.Run(); err != nil {
		panic(err)
	}
	defer client.Conn.Close()

	logger.Config(fmt.Sprintf("Client is connected to the server on port %d...", clientCfg.ServerPort))

	client.Dispatch()
}
