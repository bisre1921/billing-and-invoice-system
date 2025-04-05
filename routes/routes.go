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

func SetupCompanyRoutes(router *gin.RouterGroup) {
	company := router.Group("/company")
	{
		company.POST("/create", controllers.CreateCompany)
	}
}

func SetupEmployeeRoutes(router *gin.RouterGroup) {
	employee := router.Group("/employee")
	{
		employee.POST("/add", controllers.AddEmployee)
		employee.DELETE("/delete/:id", controllers.DeleteEmployee)
	}
}
