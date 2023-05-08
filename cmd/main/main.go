package main

import (
	"biocad/internal/app"
	"biocad/internal/config"
	"log"
)

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error in config parsing: %s\n", err)
	}

	log.Println("Application started")

	app.Run(cfg)
}
