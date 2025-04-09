package routes

import (
	"github.com/bisre1921/billing-and-invoice-system/controllers"
	"github.com/bisre1921/billing-and-invoice-system/middleware"
	"github.com/gin-gonic/gin"
)

func SetupAllRoutes(router *gin.RouterGroup) {
	_SetupAuthRoutes(router)
	_SetupCompanyRoutes(router)
	_SetupUserRoutes(router)
	_SetupEmployeeRoutes(router)
	_SetupCustomerRoutes(router)
}

func _SetupAuthRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register/user", controllers.RegisterUser)
		auth.POST("/login", controllers.Login)
	}
}

func _SetupCompanyRoutes(router *gin.RouterGroup) {
	company := router.Group("/company")
	company.Use(middleware.AuthMiddleware())
	{
		company.POST("/create", controllers.CreateCompany)
	}
}

func _SetupUserRoutes(router *gin.RouterGroup) {
	user := router.Group("/user")
	user.Use(middleware.AuthMiddleware())
	{
		user.PATCH("/update/:id", controllers.UpdateUser)
	}
}

func _SetupEmployeeRoutes(router *gin.RouterGroup) {
	employee := router.Group("/employee")
	{
		employee.POST("/add", controllers.AddEmployee)
		employee.DELETE("/delete/:id", controllers.DeleteEmployee)
		employee.GET("/all", controllers.GetAllEmployees)
	}
}

func _SetupCustomerRoutes(router *gin.RouterGroup) {
	customer := router.Group("/customer")
	{
		customer.POST("/register", controllers.RegisterCustomer)
		customer.PUT("/update/:id", controllers.UpdateCustomer)
		customer.DELETE("/delete/:id", controllers.DeleteCustomer)
	}
}
