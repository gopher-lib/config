package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gopher-lib/config"
)

func main() {
	var cfg struct {
		Port uint
	}
	err := config.LoadFile("configs/config.yaml", &cfg)
	if err != nil {
		log.Fatal()
	}

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Listening on %s...", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
