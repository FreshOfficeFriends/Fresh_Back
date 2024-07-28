package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/FreshOfficeFriends/SSO/pkg/logger"

	"github.com/FreshOfficeFriends/SSO/internal/service/auth"
)

type Handler struct {
	usersService *auth.Auth
}

func NewHandler(auth *auth.Auth) *Handler {
	return &Handler{usersService: auth}
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(Logger())

	api := r.Group("/sso")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/check-unique-email", h.uniqueEmail)
			auth.POST("/sign-up", h.signUp)              //done
			auth.POST("sign-in", h.signIn)               //done
			auth.POST("/refresh-tokens", h.refreshToken) //done

			auth.POST("/pass-recovery", h.recoverPassword)
			auth.GET("/reset-password/:uuid", h.resetPassword)
			auth.GET("/recover-password", h.r)

			auth.GET("/confirm-email/:hashEmail", h.saveUser)
		}

		//internal
		api.POST("/validate-tokens", h.validateToken)

	}
	return r
}

func (h *Handler) resetPassword(ctx *gin.Context) {
	uuid, ok := ctx.Params.Get("uuid")
	if !ok {
		ctx.JSON(http.StatusBadRequest, Response{Status: fail, Error: &ErrorDetails{"", "wrong link"}})
		return
	}

	ctx.SetCookie("uuid", uuid, 3600, "/", "localhost", false, true)

	ctx.Redirect(http.StatusPermanentRedirect, "http://localhost:8080/sso/auth/recover-password")
	//ctx.JSON(http.StatusOK, Response{
	//	Status: success,
	//	Data:   "test",
	//})
}

// поход в БД, где меняю пароль
func (h *Handler) r(ctx *gin.Context) {
	p := new(Pass)

	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Status: fail, Error: &ErrorDetails{"", err.Error()}})
		return
	}

	if err := p.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Status: fail, Error: &ErrorDetails{"", err.Error()}})
		return
	}

	uuid, err := ctx.Cookie("uuid")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Status: fail, Error: &ErrorDetails{"", "wrong link"}})
		return
	}

	//strUUID, _ := uuid.(string)
	logger.Debug(uuid)

	if err := h.usersService.ChangePass(uuid, p.Password); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Status: fail, Error: &ErrorDetails{"", err.Error()}})
		return
	}

	ctx.JSON(http.StatusOK, "all ok, pass changed")
}
