package config

import (
	"log"

	"github.com/joho/godotenv"
)

func loadEnvironmentVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type AppConfig struct {
	Logger LoggerConfig
	API    APIConfig
	Azure  AzureConfig
	Gin    GinConfig
}

func GetAppConfig(validate bool) AppConfig {
	loadEnvironmentVariables()
	LOGGER_CONF := assembleLoggerConfig()
	API_CONF := assembleAPIConfig()
	AZURE_CONF := assembleAzureConfig()
	GIN_CONF := assembleGinConfig()
	APP_CONF := AppConfig{
		Logger: LOGGER_CONF,
		API:    API_CONF,
		Azure:  AZURE_CONF,
		Gin:    GIN_CONF,
	}
	return APP_CONF
}
