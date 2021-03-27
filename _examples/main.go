package main

import (
	"log"
	"os"

	"github.com/gopher-lib/config"
)

type appconfig struct {
	Port string
}

func main() {
	os.Setenv("PORT", "")
	var conf appconfig
	err := config.LoadFile(&conf, "./configs/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
}
