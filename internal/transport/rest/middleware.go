package rest

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/FreshOfficeFriends/SSO/pkg/logger"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger.Info("router middleware",
			zap.String("method", ctx.Request.Method),
			zap.String("URI", ctx.Request.RequestURI),
			zap.String("remote address", ctx.Request.RemoteAddr))
	}
}
