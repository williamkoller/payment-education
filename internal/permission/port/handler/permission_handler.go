package port_permission_handler

import "github.com/gin-gonic/gin"

type PermissionHandler interface {
	CreatePermission(c *gin.Context)
	FindAllPermission(c *gin.Context)
	FindPermissionByUserID(c *gin.Context)
	UpdatePermission(c *gin.Context)
	DeletePermission(c *gin.Context)
	FindPermissionById(c *gin.Context)
}
