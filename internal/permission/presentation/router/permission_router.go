package permission_router

import (
	"time"

	"github.com/gin-gonic/gin"
	infra_cryptography "github.com/williamkoller/system-education/internal/auth/infra/cryptography"
	auth_middleware "github.com/williamkoller/system-education/internal/auth/presentation/middleware"
	permission_usecase "github.com/williamkoller/system-education/internal/permission/application/usecase"
	permission_repository "github.com/williamkoller/system-education/internal/permission/infra/db/repository"
	permission_handler "github.com/williamkoller/system-education/internal/permission/presentation/handler"
	permission_middleware "github.com/williamkoller/system-education/internal/permission/presentation/middleware"
	"gorm.io/gorm"
)

func PermissionRouter(e *gin.Engine, db *gorm.DB, secret string, expiresIn time.Duration) {
	repo := permission_repository.NewPermissionGormRepository(db)
	jwt := infra_cryptography.NewJWTTokenManager(secret, expiresIn)

	usecase := permission_usecase.NewPermissionUsecase(repo)
	handler := permission_handler.NewPermissionHandler(usecase)
	middleware := permission_middleware.NewPermissionMiddleware()

	p := e.Group("/permissions")
	{
		p.POST("", handler.CreatePermission)
		p.GET("", handler.FindAllPermission)
		p.GET("/user/:user_id", auth_middleware.AuthMiddleware(jwt),
			middleware.ModuleAccessMiddleware([]string{"permissions"}, []string{"read"}), handler.FindPermissionByUserID)
		p.PUT("/:id", auth_middleware.AuthMiddleware(jwt),
			middleware.ModuleAccessMiddleware([]string{"permissions"}, []string{"update"}), handler.UpdatePermission)
		p.DELETE("/:id", auth_middleware.AuthMiddleware(jwt),
			middleware.ModuleAccessMiddleware([]string{"permissions"}, []string{"delete"}), handler.DeletePermission)
		p.GET("/:id", auth_middleware.AuthMiddleware(jwt),
			middleware.ModuleAccessMiddleware([]string{"permissions"}, []string{"read"}), handler.FindPermissionById)
	}
}
