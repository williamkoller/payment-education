package user_router

import (
	"github.com/gin-gonic/gin"
	user_usecase "github.com/williamkoller/system-education/internal/user/application/usecase"
	"github.com/williamkoller/system-education/internal/user/infra/cryptography"
	user_repository "github.com/williamkoller/system-education/internal/user/infra/db/repository"
	user_handler "github.com/williamkoller/system-education/internal/user/presentation/handler"
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
		users.GET(":id", userHandler.FindByID)
		users.PUT(":id", userHandler.Update)
		users.DELETE(":id", userHandler.Delete)
	}
}
