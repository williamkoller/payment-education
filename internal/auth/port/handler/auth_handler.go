package port_auth_handler

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	Login(c *gin.Context)
	Profile( c *gin.Context)
}
