package routes

import (
	"b-resto/controllers"

	"github.com/gin-gonic/gin"
)

func SetupCompanyRoutes(router *gin.RouterGroup) {
	companies := router.Group("/companies")
	{
		companies.GET("", controllers.GetCompanies)
		companies.GET("/:id", controllers.GetCompany)
		companies.POST("", controllers.CreateCompany)
		companies.PUT("/:id", controllers.UpdateCompany)
		companies.DELETE("/:id", controllers.DeleteCompany)
	}
}
