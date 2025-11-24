package permission_router

import (
	"github.com/gin-gonic/gin"
	permission_usecase "github.com/williamkoller/system-education/internal/permission/application/usecase"
	permission_repository "github.com/williamkoller/system-education/internal/permission/infra/db/repository"
	permission_handler "github.com/williamkoller/system-education/internal/permission/presentation/handler"
	"gorm.io/gorm"
)

func PermissionRouter(e *gin.Engine, db *gorm.DB) {
	repo := permission_repository.NewPermissionGormRepository(db)

	usecase := permission_usecase.NewPermissionUsecase(repo)
	handler := permission_handler.NewPermissionHandler(usecase)
	p := e.Group("/permissions")
	{
		p.POST("", handler.CreatePermission)
		p.GET("", handler.FindAllPermission)
		p.GET("/user/:user_id", handler.FindPermissionByUserID)
		p.PUT("/:id", handler.UpdatePermission)
		p.DELETE("/:id", handler.DeletePermission)
		p.GET("/:id", handler.FindPermissionById)
	}
}
