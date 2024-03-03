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
			auth.POST("/check-unique-email", h.uniqueEmail)
			auth.POST("/sign-up", h.signUp)
			auth.GET("/confirm-email/:hashEmail", h.saveUser)
			auth.POST("sign-in", h.signIn)
			auth.POST("/refresh-tokens", h.refreshToken)
		}
		authMiddleware := api.Group("/middleware")
		{
			authMiddleware.Use(h.authMid())
			authMiddleware.GET("/test", h.test)
		}
	}
	return r
}

func (h *Handler) test(c *gin.Context) {
	id, _ := c.Get("userId")
	c.JSON(200, gin.H{"userId": id})
}

//дореализовать sign-in
//
//вынести валидацию токена в отдельный ендпоинт
//настроить проверку токена в nginx
