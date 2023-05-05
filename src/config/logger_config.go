package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type LogLevelEnum string

const (
	DEBUG   LogLevelEnum = "DEBUG"
	INFO    LogLevelEnum = "INFO"
	WARNING LogLevelEnum = "WARNING"
	ERROR   LogLevelEnum = "ERROR"
)

type LogLevel struct {
	LogLevelEnum
}

func (l *LogLevel) FromStr(levelStr string) error {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		l.LogLevelEnum = DEBUG
	case "INFO":
		l.LogLevelEnum = INFO
	case "WARNING":
		l.LogLevelEnum = WARNING
	case "ERROR":
		l.LogLevelEnum = ERROR
	default:
		allowedOptions := "['DEBUG', 'INFO', 'WARNING', 'ERROR']"
		msg := "failed to determine provided LogLevel from string;"
		errMsg := fmt.Sprintf("%s acceptable values are %s", msg, allowedOptions)
		return errors.New(errMsg)
	}
	return nil
}

type LoggerConfig struct {
	LogLevel LogLevel `validate:"required"`
	FilePath string   `validate:"required,filepath"`
}

func assembleLoggerConfig() LoggerConfig {
	// Collect Log Level
	LOG_LEVEL := LogLevel{}
	err := LOG_LEVEL.FromStr(os.Getenv(("LOGGER.LOG_LEVEL")))
	if err != nil {
		log.Fatal(err.Error())
	}
	return LoggerConfig{
		LogLevel: LOG_LEVEL,
		FilePath: os.Getenv("LOGGER.FILE_PATH"),
	}
}
