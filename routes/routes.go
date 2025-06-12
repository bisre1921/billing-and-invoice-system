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
	_SetupItemRoutes(router)
	_SetupInvoiceRoutes(router)
	_SetupReportRoutes(router)
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
		company.GET("/:id", controllers.GetCompany)
		company.GET("user/:user_id", controllers.CheckCompanyForUser)
	}
}

func _SetupUserRoutes(router *gin.RouterGroup) {
	user := router.Group("/user")
	user.Use(middleware.AuthMiddleware())
	{
		user.PATCH("/update/:id", controllers.UpdateUser)
		user.GET("/:id", controllers.GetUser)
	}
}

func _SetupEmployeeRoutes(router *gin.RouterGroup) {
	employee := router.Group("/employee")
	{
		employee.POST("/add", controllers.AddEmployee)
		employee.DELETE("/delete/:id", controllers.DeleteEmployee)
		employee.GET("/all", controllers.GetAllEmployees)
		employee.GET("/:id", controllers.GetEmployee)
		employee.PUT("/update/:id", controllers.UpdateEmployee)
	}
}

func _SetupCustomerRoutes(router *gin.RouterGroup) {
	customer := router.Group("/customer")
	{
		customer.POST("/register", controllers.RegisterCustomer)
		customer.PUT("/update/:id", controllers.UpdateCustomer)
		customer.DELETE("/delete/:id", controllers.DeleteCustomer)
		customer.GET("/all", controllers.ListCustomers)
		customer.GET("/:id", controllers.GetCustomer)
	}
}

func _SetupItemRoutes(router *gin.RouterGroup) {
	item := router.Group("/item")
	{
		item.POST("/add", controllers.AddItem)
		item.PUT("/update/:id", controllers.UpdateItem)
		item.DELETE("/delete/:id", controllers.DeleteItem)
		item.GET("/all", controllers.ListItems)
		item.GET("/company/:company_id", controllers.GetItemsByCompanyID)
		item.GET("/:id", controllers.GetItem)
		item.POST("/import", controllers.ImportItems)
	}
}

func _SetupInvoiceRoutes(router *gin.RouterGroup) {
	invoice := router.Group("/invoice")
	{
		invoice.POST("/generate", controllers.GenerateInvoice)
		invoice.GET("/:id", controllers.GetInvoice)
		invoice.GET("/companies/:company_id", controllers.GetInvoicesByCompanyID)
		invoice.POST("/send/:id", controllers.SendInvoice)
		invoice.GET("/download/:id", controllers.DownloadInvoice)
		invoice.PUT("/mark-as-paid/:id", controllers.MarkInvoiceAsPaid)
	}
}

func _SetupReportRoutes(router *gin.RouterGroup) {
	report := router.Group("/report")
	{
		report.POST("/sales", controllers.GetSalesReport)
	}
}
