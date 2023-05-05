package config

import (
	"log"
	"os"
	"strconv"
)

type APIConfig struct {
	Port int
}

func assembleAPIConfig() APIConfig {
	PORT, err := strconv.Atoi(os.Getenv("API.PORT"))
	if err != nil {
		log.Fatal(err.Error())
	}
	return APIConfig{
		Port: PORT,
	}
}
