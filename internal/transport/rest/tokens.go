package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/FreshOfficeFriends/SSO/internal/domain"
	"github.com/FreshOfficeFriends/SSO/pkg/logger"
)

func (h *Handler) refreshToken(ctx *gin.Context) {
	cookie, err := ctx.Cookie("refresh-token")
	if err != nil {
		logger.Error(err.Error())
		ctx.Status(http.StatusBadRequest)
		return
	}

	logger.Debug(cookie)

	accessToken, refreshToken, err := h.usersService.RefreshTokens(cookie)
	if err != nil {
		if err == domain.TokenExpired {
			ctx.JSON(http.StatusUnauthorized, Response{Status: fail, Error: &ErrorDetails{"",
				"token expired, try to log in again"}})
			return
		}
		logger.Error(err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	}
	logger.Debug(fmt.Sprintf("access=%s, refresh=%s", accessToken, refreshToken))

	ctx.SetCookie("refresh-token", refreshToken, 5184000, "localhost", "/", false, true)
	ctx.JSON(http.StatusOK, Response{Status: success, Data: map[string]string{"accessToken": accessToken}})
}

func (h *Handler) validateToken(ctx *gin.Context) {
	accessToken, err := getTokenFromRequest(ctx.Request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, Response{Status: fail, Error: &ErrorDetails{"", err.Error()}})
		logger.Error(err.Error())
		return
	}

	userId, err := h.usersService.ParseToken(accessToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, Response{Status: fail, Error: &ErrorDetails{"", "invalid token"}})
		logger.Error(err.Error())
		return
	}

	ctx.JSON(http.StatusOK, Response{Status: success, Data: map[string]int{"userId": userId}})
}

func getTokenFromRequest(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if len(token) == 0 {
		return "", domain.EmptyTokenHeader
	}

	headerAuth := strings.Split(token, " ")
	if headerAuth[0] != "Bearer" || len(headerAuth) != 2 {
		logger.Debug("split token", zap.String("", headerAuth[0]+headerAuth[1]))
		return "", domain.InvalidHeaderSignature
	}
	//todo передай бирер и пустой токен

	return headerAuth[1], nil
}
