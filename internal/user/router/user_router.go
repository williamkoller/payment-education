package user_router

import (
	"github.com/gin-gonic/gin"
	"github.com/williamkoller/system-education/internal/user/cryptography"
	user_handler "github.com/williamkoller/system-education/internal/user/handler"
	user_repository "github.com/williamkoller/system-education/internal/user/repository"
	user_usecase "github.com/williamkoller/system-education/internal/user/usecase"
	"gorm.io/gorm"
)

func UserRouter(e *gin.Engine, db *gorm.DB) {
	crypto := cryptography.NewBcryptHasher(12)
	userRepo := user_repository.NewUserGormRepository(db)
	userUsecase := user_usecase.NewUserUsecase(userRepo, crypto)
	userHandler := user_handler.NewUserHandler(userUsecase)
	users := e.Group("/users")
	{
		users.POST("", userHandler.CreateUser)
		users.GET("", userHandler.FindAllUsers)
	}
}
