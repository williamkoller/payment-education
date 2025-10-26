package user_router

import "github.com/gin-gonic/gin"

func UserRouter(e *gin.Engine) {
	users := e.Group("/users")
	{
		users.GET("/", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"message": "ok"})
		})
	}
}