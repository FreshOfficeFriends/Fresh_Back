package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/FreshOfficeFriends/SSO/internal/domain"
	"github.com/FreshOfficeFriends/SSO/pkg/logger"
)

func (h *Handler) uniqueEmail(ctx *gin.Context) {
	email := new(Email)

	_ = ctx.BindJSON(&email)

	logger.Debug(email.Email)

	if err := email.Validate(); err != nil {
		rsp := Response{
			fail,
			nil,
			&ErrorDetails{"", "empty 'email' field or invalid 'email' format"},
		}
		ctx.JSON(http.StatusBadRequest, rsp)
		logger.Error("empty email or invalid format")
		return
	}

	err := h.usersService.UniqueEmail(email.Email)
	if err != nil {
		rsp := Response{
			fail,
			nil,
			&ErrorDetails{"", err.Error()},
		}
		ctx.AbortWithStatusJSON(http.StatusConflict, rsp)
		logger.Error(err.Error())
		return
	}

	rsp := Response{
		success,
		"email is unique",
		nil,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (h *Handler) signUp(ctx *gin.Context) {
	user := new(domain.SignUp)

	_ = ctx.ShouldBind(&user)
	logger.Debug(fmt.Sprintf("signUp data - %+v", user))

	if err := user.Validate(); err != nil {
		rsp := Response{
			fail,
			nil,
			&ErrorDetails{"", fmt.Sprintf("data has not been validated, error - %s", err.Error())},
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, rsp)
		logger.Error(err.Error())
		return
	}

	if err := h.usersService.SignUp(user); err != nil {
		rsp := Response{
			fail,
			nil,
			&ErrorDetails{"", "registration error (on server)"},
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, rsp)
		logger.Error(err.Error())
		return
	}

	ctx.JSON(http.StatusOK, Response{Status: success, Data: "An email has been sent to the user to confirm registration"})
}

func (h *Handler) saveUser(ctx *gin.Context) {
	val, ok := ctx.Params.Get("hashEmail")
	if !ok {
		ctx.JSON(http.StatusBadRequest, Response{Status: fail, Error: &ErrorDetails{"", "wrong link"}})
		return
	}

	if err := h.usersService.SaveUser(val); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Status: fail, Error: &ErrorDetails{"", "server error, try registering again"}})
		logger.Error(err.Error())
		return
	}

	defer ctx.Redirect(http.StatusPermanentRedirect, "https://apple.com")
}

func (h *Handler) signIn(ctx *gin.Context) {
	user := new(domain.SignIn)

	if err := ctx.ShouldBind(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Error(err.Error())
		return
	}

	accessToken, refreshToken, err := h.usersService.SignIn(user)

	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logger.Error(err.Error())
		return
	}

	ctx.SetCookie("refresh-token", refreshToken, 3600, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, Response{Status: success, Data: map[string]string{"accessToken": accessToken}})
}

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

	ctx.SetCookie("refresh-token", refreshToken, 3600, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, Response{Status: success, Data: map[string]string{"accessToken": accessToken}})
}
