package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (env *AppEnvironment) MiddlewareInboundRequestLog(context *gin.Context) {
	env.Logger.WithFields(
		logrus.Fields{
			"http_method": context.Request.Method,
			"client_ip":   context.ClientIP(),
			"host":        context.Request.Host,
			"headers":     context.Request.Header,
		},
	).Info("Inbound Request")
	context.Next()
}
