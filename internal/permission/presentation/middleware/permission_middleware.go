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
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permissions not found in token"})
			return
		}

		modulesSlice, ok := modulesInterface.([]interface{})
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid permissions format"})
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
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied to required modules"})
			return
		}

		// Validate actions (if required)
		if len(requiredActions) > 0 {
			actionsInterface, ok := c.Get("actions")
			if !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Actions not found in token"})
				return
			}

			actionsSlice, ok := actionsInterface.([]interface{})
			if !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid actions format"})
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
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied to required actions"})
				return
			}
		}

		c.Next()
	}
}
