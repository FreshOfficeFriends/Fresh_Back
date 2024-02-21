package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/FreshOfficeFriends/SSO/internal/domain"
	"github.com/FreshOfficeFriends/SSO/pkg/logger"
)

func (h *Handler) uniqueEmail(ctx *gin.Context) {
	var emailData map[string]string

	if err := ctx.BindJSON(&emailData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Error(err.Error())
		return
	}

	email, ok := emailData["email"]
	if !ok || email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "'email' field is required"})
		logger.Error("'email' field error")
		return
	}

	err := h.usersService.UniqueEmail(email)
	if err != nil {
		ctx.Status(http.StatusConflict)
		logger.Error(err.Error())
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *Handler) signUp(ctx *gin.Context) {
	user := new(domain.SignUp)

	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Error(err.Error())
		return
	}

	logger.Debug(fmt.Sprintf("%+v", user))

	//if err := user.Validate(); err != nil {
	//	ctx.Status(http.StatusBadRequest)
	//	logger.Error(err.Error())
	//	return
	//}

	if err := h.usersService.SignUp(user); err != nil {
		ctx.Status(http.StatusInternalServerError)
		logger.Error(err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) saveUser(ctx *gin.Context) {
	val, ok := ctx.Params.Get("hashEmail")
	if !ok {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if err := h.usersService.SaveUser(val); err != nil {
		ctx.Status(http.StatusInternalServerError)
		logger.Error(err.Error())
		return
	}

	defer ctx.Redirect(http.StatusPermanentRedirect, "https://apple.com")
}
