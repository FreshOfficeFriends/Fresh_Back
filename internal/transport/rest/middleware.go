package rest

import (
	"errors"
	"net/http"
	"strings"

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

func (h *Handler) authMid() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, err := getTokenFromRequest(ctx.Request)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			logger.Error(err.Error())
			return
		}

		userId, err := h.usersService.ParseToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			logger.Error(err.Error())
			return
		}

		ctx.Set("userId", userId)
		ctx.Next()
	}
}
func getTokenFromRequest(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if len(token) == 0 {
		return "", errors.New("empty token, add auth header")
	}

	headerAuth := strings.Split(token, " ")
	if headerAuth[0] != "Bearer" || len(headerAuth) != 2 {
		logger.Debug("split token", zap.String("", headerAuth[0]+headerAuth[1]))
		return "", errors.New("invalid auth header, try 'Bearer <token>'")
	}
	//todo передай бирер и пустой токен

	return headerAuth[1], nil
}
