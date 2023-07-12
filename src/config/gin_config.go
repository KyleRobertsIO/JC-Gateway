package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
)

type GinConfig struct {
	Port int
	Mode string
}

func (gc *GinConfig) validateGinMode(mode string) error {
	switch strings.ToLower(mode) {
	case "debug":
		gc.Mode = "debug"
	case "release":
		gc.Mode = "release"
	default:
		allowedOptions := "['debug', 'release']"
		msg := "failed to determine provided GinMode from string;"
		errMsg := fmt.Sprintf("%s acceptable values are %s", msg, allowedOptions)
		return errors.New(errMsg)
	}
	return nil
}

func assembleGinConfig() GinConfig {
	GIN_CONFIG := GinConfig{}
	modeErr := GIN_CONFIG.validateGinMode(os.Getenv(("GIN_MODE")))
	if modeErr != nil {
		log.Fatal(modeErr.Error())
	}
	PORT, err := strconv.Atoi(os.Getenv("GIN_PORT"))
	if err != nil {
		log.Fatal(err.Error())
	}
	GIN_CONFIG.Port = PORT
	return GIN_CONFIG
}
