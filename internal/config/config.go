package config

import (
	"github.com/caarlos0/env/v7"
	"log"
)

type Config struct {
	AppPort     string `env:"APP_PORT" envDefault:"9000"`
	PostgresUrl string `env:"BIOCAD_DB_URL"`
	DirAddress  string `env:"DIR_ADDRESS"`
	PDFSaveAddr string `env:"PDF_DIR_ADDRESS"`
	FontAddress string `env:"FONT_DIR_ADDRESS"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		log.Println("Error in parsing env")
		return nil, err
	}
	//pwd, err := os.Getwd()
	//if err != nil {
	//	log.Println("Error in getting pwd")
	//	return nil, err
	//}
	//cfg.DirAddress = pwd + "/resources/testData"
	//cfg.PDFSaveAddr = pwd + "/resources/pdf"

	//for testing purposes
	//cfg.PostgresUrl = "postgresql://biocad:password@localhost:5432/biocad"

	return cfg, nil
}
