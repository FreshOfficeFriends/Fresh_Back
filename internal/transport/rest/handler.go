package rest

import (
	"github.com/gin-gonic/gin"

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
			auth.POST("/unique-email", h.uniqueEmail)
			auth.POST("/sign-up", h.signUp)
			auth.GET("/confirm-email/:hashEmail", h.saveUser)
		}
	}
	return r
}
