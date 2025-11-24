package user_router

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	infra_cryptography "github.com/williamkoller/system-education/internal/auth/infra/cryptography"
	auth_middleware "github.com/williamkoller/system-education/internal/auth/presentation/middleware"
	permission_middleware "github.com/williamkoller/system-education/internal/permission/presentation/middleware"
	user_usecase "github.com/williamkoller/system-education/internal/user/application/usecase"
	user_event "github.com/williamkoller/system-education/internal/user/domain/event"
	user_cryptography "github.com/williamkoller/system-education/internal/user/infra/cryptography"
	user_repository "github.com/williamkoller/system-education/internal/user/infra/db/repository"
	infra_email "github.com/williamkoller/system-education/internal/user/infra/email"
	user_handler "github.com/williamkoller/system-education/internal/user/presentation/handler"
	shared_event "github.com/williamkoller/system-education/shared/domain/event"
	"github.com/williamkoller/system-education/shared/infra/email"
	"gorm.io/gorm"
)

func UserRouter(e *gin.Engine, db *gorm.DB, apiKey string, fromAddress string, secret string, expiresIn time.Duration) {
	crypto := user_cryptography.NewBcryptHasher(12)
	userRepo := user_repository.NewUserGormRepository(db)
	event := shared_event.NewDispatcher()
	jwt := infra_cryptography.NewJWTTokenManager(secret, expiresIn)
	middleware := permission_middleware.NewPermissionMiddleware()

	client := email.NewResendClient(apiKey, fromAddress)
	notifier := infra_email.NewResendEmailNotifier(client)
	event.Register("user.created", func(e interface{}) {
		evt, ok := e.(*user_event.UserCreatedEvent)
		if !ok {
			log.Printf("Evento inesperado: %+v", e)
			return
		}

		if err := notifier.SendWelcomeEmail(evt.Name, evt.Email); err != nil {
			log.Printf("Falha ao enviar e‑mail de boas‑vindas: %v", err)
		} else {
			log.Printf("E‑mail de boas‑vindas enviado para: %s", evt.Email)
		}
	})

	userUsecase := user_usecase.NewUserUsecase(userRepo, crypto, event)
	userHandler := user_handler.NewUserHandler(userUsecase)

	users := e.Group("/users")
	{
		users.POST("", userHandler.CreateUser)
		users.GET("", userHandler.FindAllUsers)
		users.GET(":id",
			auth_middleware.AuthMiddleware(jwt),
			middleware.ModuleAccessMiddleware([]string{"users"}),
			userHandler.FindByID,
		)
		users.PUT(":id",
			auth_middleware.AuthMiddleware(jwt),
			middleware.ModuleAccessMiddleware([]string{"users"}),
			userHandler.Update,
		)
		users.DELETE(":id",
			auth_middleware.AuthMiddleware(jwt),
			middleware.ModuleAccessMiddleware([]string{"users"}),
			userHandler.Delete,
		)
	}
}
