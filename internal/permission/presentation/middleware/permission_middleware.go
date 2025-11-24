package permission_middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	port_permission_middleware "github.com/williamkoller/system-education/internal/permission/port/middleware"
)

var _ port_permission_middleware.PermissionMiddleware = &PermissionMiddleware{}

type PermissionMiddleware struct{}

func NewPermissionMiddleware() *PermissionMiddleware {
	return &PermissionMiddleware{}
}

func (m *PermissionMiddleware) ModuleAccessMiddleware(requiredModules []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		modulesInterface, ok := c.Get("modules")
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permissões não encontradas no token"})
			return
		}

		modulesSlice, ok := modulesInterface.([]interface{})
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Formato de permissões inválido"})
			return
		}

		userModules := make([]string, 0, len(modulesSlice))
		for _, mod := range modulesSlice {
			if modStr, ok := mod.(string); ok {
				userModules = append(userModules, modStr)
			}
		}

		for _, required := range requiredModules {
			for _, userMod := range userModules {
				if userMod == required {
					c.Next()
					return
				}
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Acesso negado aos módulos exigidos"})
	}
}
