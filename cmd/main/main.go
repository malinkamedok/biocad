package main

import (
	"biocad/internal/app"
	"biocad/internal/configs"
	"log"
)

func main() {

	cfg, err := configs.NewConfig()
	if err != nil {
		log.Fatalf("Error in config parsing: %s\n", err)
	}

	log.Println("Application started")

	app.Run(cfg)
}
