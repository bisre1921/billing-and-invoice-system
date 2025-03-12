package routes

import (
	"github.com/bisre1921/billing-and-invoice-system/controllers"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register/user", controllers.RegisterUser)
		auth.POST("/login", controllers.Login)
	}
}
