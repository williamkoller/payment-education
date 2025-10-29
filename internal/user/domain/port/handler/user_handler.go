package port_user_handler

import "github.com/gin-gonic/gin"

type UserHandler interface {
	CreateUser(c *gin.Context)
}
