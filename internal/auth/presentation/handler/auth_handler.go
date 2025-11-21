package auth_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	port_auth_handler "github.com/williamkoller/system-education/internal/auth/port/handler"
	port_auth_usecase "github.com/williamkoller/system-education/internal/auth/port/usecase"
	auth_dtos "github.com/williamkoller/system-education/internal/auth/presentation/dtos"
)

type AuthHandler struct {
	usecase port_auth_usecase.AuthUsecase
}

func NewAuthHandler(usecase port_auth_usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{usecase: usecase}
}

var _ port_auth_handler.AuthHandler = &AuthHandler{}

func (h *AuthHandler) Login(c *gin.Context) {
	var input auth_dtos.AuthDto

	if err := c.ShouldBindJSON(&input); err != nil {
		c.Status(http.StatusBadRequest)
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	token, err := h.usecase.Login(input.Email, input.Password)

	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userEmail, _ := c.Get("userEmail")
	email, err := h.usecase.Profile(userEmail.(string))

	if err != nil {
		c.Status(http.StatusBadRequest)
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to your profile",
		"email":   email,
	})
}
