package main

import (
	"fmt"

	"github.com/charmingruby/gdp/config"
	"github.com/charmingruby/gdp/internal/network/udp"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	serverCfg := cfg.Server

	server, err := udp.NewServer(udp.ServerInput{
		Port: serverCfg.Port,
		Threshold: udp.CongestionThreshold{
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

	fmt.Println("Server is listening on port", serverCfg.Port)
	fmt.Println("Server is ready to receive data...")

	server.Read()
}
