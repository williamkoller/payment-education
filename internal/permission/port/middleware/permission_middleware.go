package port_permission_middleware

import (
	"github.com/gin-gonic/gin"
)

type PermissionMiddleware interface {
	ModuleAccessMiddleware(requiredModules []string, requiredActions []string) gin.HandlerFunc
}
