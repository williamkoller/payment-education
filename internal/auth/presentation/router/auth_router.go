package auth_router

import (
	"time"

	"github.com/gin-gonic/gin"
	auth_usecase "github.com/williamkoller/system-education/internal/auth/application/usecase"
	infra_cryptography "github.com/williamkoller/system-education/internal/auth/infra/cryptography"
	auth_handler "github.com/williamkoller/system-education/internal/auth/presentation/handler"
	auth_middleware "github.com/williamkoller/system-education/internal/auth/presentation/middleware"
	permission_repository "github.com/williamkoller/system-education/internal/permission/infra/db/repository"
	user_cryptography "github.com/williamkoller/system-education/internal/user/infra/cryptography"
	user_repository "github.com/williamkoller/system-education/internal/user/infra/db/repository"
	"gorm.io/gorm"
)

func AuthRouter(r *gin.Engine, db *gorm.DB, secret string, expiresIn time.Duration) {
	repository := user_repository.NewUserGormRepository(db)
	permissionRepo := permission_repository.NewPermissionGormRepository(db)
	jwt := infra_cryptography.NewJWTTokenManager(secret, expiresIn)
	crypto := user_cryptography.NewBcryptHasher(12)

	usecase := auth_usecase.NewAuthUsecase(repository, permissionRepo, jwt, crypto)
	handler := auth_handler.NewAuthHandler(usecase)
	auth := r.Group("auth")
	{
		auth.POST("login", handler.Login)
		auth.POST("profile", auth_middleware.AuthMiddleware(jwt), handler.Profile)
	}
}
