package permission_middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewPermissionMiddleware(t *testing.T) {
	middleware := NewPermissionMiddleware()

	assert.NotNil(t, middleware)
	assert.IsType(t, &PermissionMiddleware{}, middleware)
}

func TestModuleAccessMiddleware_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"admin", "user"}

	// Create a test router
	router := gin.New()

	// Set modules in context BEFORE the route (simulating AuthMiddleware)
	router.Use(func(c *gin.Context) {
		c.Set("modules", []interface{}{"admin", "reports"})
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, []string{}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestModuleAccessMiddleware_NoModulesInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"admin"}

	router := gin.New()
	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, []string{}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Permissões não encontradas no token")
}

func TestModuleAccessMiddleware_InvalidModulesFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"admin"}

	router := gin.New()

	// Set invalid format BEFORE the route
	router.Use(func(c *gin.Context) {
		c.Set("modules", "invalid-format")
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, []string{}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Formato de permissões inválido")
}

func TestModuleAccessMiddleware_UserHasRequiredModule(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"admin", "user"}

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("modules", []interface{}{"user", "reports"})
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, []string{}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestModuleAccessMiddleware_UserDoesNotHaveRequiredModule(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"admin", "superuser"}

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("modules", []interface{}{"user", "reports"})
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, []string{}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Acesso negado aos módulos exigidos")
}

func TestModuleAccessMiddleware_EmptyModulesList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"admin"}

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("modules", []interface{}{})
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, []string{}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Acesso negado aos módulos exigidos")
}

func TestModuleAccessMiddleware_MultipleModulesMatch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"admin", "user", "reports"}

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("modules", []interface{}{"admin", "user", "reports", "analytics"})
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, []string{}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestModuleAccessMiddleware_ModulesWithNonStringValues(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"admin"}

	router := gin.New()

	router.Use(func(c *gin.Context) {
		// Mix of string and non-string values
		c.Set("modules", []interface{}{"user", 123, "reports", nil})
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, []string{}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should fail because "admin" is not in the valid string modules
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Acesso negado aos módulos exigidos")
}

func TestModuleAccessMiddleware_NilModuleValue(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"admin"}

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("modules", nil)
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, []string{}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Formato de permissões inválido")
}

func TestModuleAccessMiddleware_SingleRequiredModuleMatch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"reports"}

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("modules", []interface{}{"reports"})
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, []string{}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestModuleAccessMiddleware_CaseSensitiveModules(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"Admin"}

	router := gin.New()

	router.Use(func(c *gin.Context) {
		// User has "admin" (lowercase), but required is "Admin" (capitalized)
		c.Set("modules", []interface{}{"admin"})
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, []string{}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should fail because module names are case-sensitive
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Acesso negado aos módulos exigidos")
}

func TestModuleAccessMiddleware_SpecialCharactersInModules(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"admin-panel", "user_management"}

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("modules", []interface{}{"admin-panel", "reports"})
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, []string{}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestModuleAccessMiddleware_EmptyRequiredModules(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{} // Empty required modules

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("modules", []interface{}{"admin", "user"})
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, []string{}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// With empty required modules, should deny access (no match found)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Acesso negado aos módulos exigidos")
}

func TestModuleAccessMiddleware_NumericModuleNames(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"module1", "module2"}

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("modules", []interface{}{"module1", "module3"})
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, []string{}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestModuleAccessMiddleware_FirstMatchWins(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"admin", "user", "reports"}

	router := gin.New()

	router.Use(func(c *gin.Context) {
		// User has "user" which is the second required module
		c.Set("modules", []interface{}{"user"})
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, []string{}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

// Action Validation Tests

func TestModuleAccessMiddleware_WithValidAction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"admin"}
	requiredActions := []string{"read", "delete"}

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("modules", []interface{}{"admin"})
		c.Set("actions", []interface{}{"read", "update"})
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, requiredActions), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestModuleAccessMiddleware_WithoutRequiredAction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"admin"}
	requiredActions := []string{"delete"}

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("modules", []interface{}{"admin"})
		c.Set("actions", []interface{}{"read", "update"})
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, requiredActions), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Acesso negado às ações exigidas")
}

func TestModuleAccessMiddleware_NoActionsInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"admin"}
	requiredActions := []string{"read"}

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("modules", []interface{}{"admin"})
		// No actions set in context
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, requiredActions), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Ações não encontradas no token")
}

func TestModuleAccessMiddleware_InvalidActionsFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"admin"}
	requiredActions := []string{"read"}

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("modules", []interface{}{"admin"})
		c.Set("actions", "invalid-format")
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, requiredActions), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Formato de ações inválido")
}

func TestModuleAccessMiddleware_EmptyActionsInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"admin"}
	requiredActions := []string{"read"}

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("modules", []interface{}{"admin"})
		c.Set("actions", []interface{}{})
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, requiredActions), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Acesso negado às ações exigidas")
}

func TestModuleAccessMiddleware_MultipleActionsOneMatches(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"admin"}
	requiredActions := []string{"read", "delete", "update"}

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("modules", []interface{}{"admin"})
		c.Set("actions", []interface{}{"delete"})
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, requiredActions), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestModuleAccessMiddleware_ActionsWithNonStringValues(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"admin"}
	requiredActions := []string{"read"}

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("modules", []interface{}{"admin"})
		// Mix of string and non-string values
		c.Set("actions", []interface{}{123, "delete", nil})
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, requiredActions), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should fail because "read" is not in the valid string actions
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Acesso negado às ações exigidas")
}

func TestModuleAccessMiddleware_BackwardCompatibility_EmptyRequiredActions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"admin"}
	requiredActions := []string{} // Empty - should skip action validation

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("modules", []interface{}{"admin"})
		// No actions set - but should still pass because requiredActions is empty
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, requiredActions), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestModuleAccessMiddleware_CaseSensitiveActions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"admin"}
	requiredActions := []string{"Read"}

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("modules", []interface{}{"admin"})
		// User has "read" (lowercase), but required is "Read" (capitalized)
		c.Set("actions", []interface{}{"read"})
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, requiredActions), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should fail because action names are case-sensitive
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Acesso negado às ações exigidas")
}

func TestModuleAccessMiddleware_AllCRUDActions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"users"}
	requiredActions := []string{"create", "read", "update", "delete"}

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("modules", []interface{}{"users"})
		c.Set("actions", []interface{}{"create", "read", "update", "delete"})
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, requiredActions), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestModuleAccessMiddleware_ModuleMatchButNoActionMatch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := NewPermissionMiddleware()
	requiredModules := []string{"admin"}
	requiredActions := []string{"delete"}

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("modules", []interface{}{"admin"}) // Module matches
		c.Set("actions", []interface{}{"read"})  // Action does not match
		c.Next()
	})

	router.GET("/test", middleware.ModuleAccessMiddleware(requiredModules, requiredActions), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should fail because action doesn't match even though module does
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Acesso negado às ações exigidas")
}
