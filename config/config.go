package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	Database
	Api
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	return &Config{
		Database: DataStore(),
		Api:      API(),
	}
}
