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

	serverCfg := cfg.Server

	server, err := udp.NewServer(udp.ServerInput{
		Port: serverCfg.Port,
	})
	if err != nil {
		panic(err)
	}

	conn, err := server.Start()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	println("UDP server is running on port", serverCfg.Port)
}
