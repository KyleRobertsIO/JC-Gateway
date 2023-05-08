package api

import "kyleroberts.io/src/config"

type AppEnvironment struct {
	Config           config.AppConfig
	AzureAccessToken string
}
