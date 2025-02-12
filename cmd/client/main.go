package main

import (
	"github.com/charmingruby/gdp/config"
	"github.com/charmingruby/gdp/internal/network/udp"
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

	println("Client is connected to the server on port", clientCfg.ServerPort)
	println("Client is ready to communicate...")

	client.Dispatch()
}
