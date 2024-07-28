package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/FreshOfficeFriends/SSO/internal/service/auth"

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
			&ErrorDetails{"", domain.BadEmail.Error()},
		}
		ctx.JSON(http.StatusBadRequest, rsp)
		logger.Error("empty email or invalid format")
		return
	}

	uniqueEmail := h.usersService.UniqueEmail(email.Email)
	if !uniqueEmail {
		rsp := Response{
			fail,
			nil,
			&ErrorDetails{"", "email is not unique"},
		}
		ctx.AbortWithStatusJSON(http.StatusConflict, rsp)
		logger.Error("email is not unique")
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

// принимаю почту, отправляю на нее ссылку для восстановления пароля
// проверка, что такая почта существует
// ссылка - http://fresh/sso/recover-password
func (h *Handler) recoverPassword(ctx *gin.Context) {
	email := new(Email)

	_ = ctx.BindJSON(&email)

	logger.Debug(email.Email)

	if err := email.Validate(); err != nil {
		rsp := Response{
			fail,
			nil,
			&ErrorDetails{"", domain.BadEmail.Error()},
		}
		ctx.JSON(http.StatusBadRequest, rsp)
		logger.Error("empty email or invalid format")
		return
	}

	uuid, err := h.usersService.UUID(email.Email)
	if err != nil {
		rsp := Response{
			fail,
			nil,
			&ErrorDetails{"", "email not exist"},
		}
		ctx.AbortWithStatusJSON(http.StatusConflict, rsp)
		logger.Error("email not exist")
		return
	}

	err = auth.SendEmailRecoverPass(email.Email, uuid)
	if err != nil {
	}
	rsp := Response{
		success,
		"email has been sent, check your inbox",
		nil,
	}
	ctx.JSON(http.StatusOK, rsp)

	//email либо отправлен, либо нет
	//отправлен = указан именно email + такой email есть в базе
}
