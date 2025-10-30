package user_router

import (
	"github.com/gin-gonic/gin"
	user_handler "github.com/williamkoller/system-education/internal/user/handler"
	user_repository "github.com/williamkoller/system-education/internal/user/repository"
	user_usecase "github.com/williamkoller/system-education/internal/user/usecase"
)

func UserRouter(e *gin.Engine) {
	userRepo := user_repository.NewUserMemoryRepository()
	userUsecase := user_usecase.NewUserUsecase(userRepo)
	userHandler := user_handler.NewUserHandler(userUsecase)
	users := e.Group("/users")
	{
		users.POST("", userHandler.CreateUser)
		users.GET("", userHandler.FindAllUsers)
	}
}
