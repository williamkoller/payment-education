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

func (m *PermissionMiddleware) ModuleAccessMiddleware(requiredModules []string, requiredActions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate modules
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

		// Check if user has at least one required module
		hasModule := false
		for _, required := range requiredModules {
			for _, userMod := range userModules {
				if userMod == required {
					hasModule = true
					break
				}
			}
			if hasModule {
				break
			}
		}

		if !hasModule {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Acesso negado aos módulos exigidos"})
			return
		}

		// Validate actions (if required)
		if len(requiredActions) > 0 {
			actionsInterface, ok := c.Get("actions")
			if !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Ações não encontradas no token"})
				return
			}

			actionsSlice, ok := actionsInterface.([]interface{})
			if !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Formato de ações inválido"})
				return
			}

			userActions := make([]string, 0, len(actionsSlice))
			for _, act := range actionsSlice {
				if actStr, ok := act.(string); ok {
					userActions = append(userActions, actStr)
				}
			}

			// Check if user has at least one required action
			hasAction := false
			for _, required := range requiredActions {
				for _, userAct := range userActions {
					if userAct == required {
						hasAction = true
						break
					}
				}
				if hasAction {
					break
				}
			}

			if !hasAction {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Acesso negado às ações exigidas"})
				return
			}
		}

		c.Next()
	}
}
