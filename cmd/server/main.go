package main

import (
	"fmt"

	"github.com/charmingruby/gdp/config"
	"github.com/charmingruby/gdp/internal/network/udp/server"
	"github.com/charmingruby/gdp/pkg/logger"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	serverCfg := cfg.Server

	fmt.Printf("Client configuration: %+v\n", serverCfg)

	server, err := server.New(server.ServerInput{
		Port: serverCfg.Port,
		Threshold: server.CongestionThreshold{
			PackageLoss: serverCfg.CongestionThreshold.PackageLoss,
		},
	})
	if err != nil {
		panic(err)
	}

	if err := server.Listen(); err != nil {
		panic(err)
	}
	defer server.Conn.Close()

	logger.Config(fmt.Sprintf("Server is listening on port %d...", serverCfg.Port))

	if err := server.Read(); err != nil {
		panic(err)
	}
}
