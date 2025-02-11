package main

import (
	"fmt"

	"github.com/charmingruby/gdp/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", cfg)
}
