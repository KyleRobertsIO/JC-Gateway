package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (env *AppEnvironment) Ping(context *gin.Context) {
	context.JSON(
		http.StatusOK,
		gin.H{"message": "pong"},
	)
}
