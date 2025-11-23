package port_permission_handler

import "github.com/gin-gonic/gin"

type PermissionHandler interface {
	CreatePermission(c *gin.Context)
	FindAllPermission(c *gin.Context)
}
