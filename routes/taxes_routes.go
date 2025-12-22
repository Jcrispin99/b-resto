package routes

import (
	"b-resto/controllers"

	"github.com/gin-gonic/gin"
)

func SetupTaxesRoutes(router *gin.RouterGroup) {
	taxes := router.Group("/taxes")
	{
		taxes.GET("", controllers.GetTaxes)
		taxes.GET("/:id", controllers.GetTax)
		taxes.POST("", controllers.CreateTax)
		taxes.PUT("/:id", controllers.UpdateTax)
		taxes.DELETE("/:id", controllers.DeleteTax)
		taxes.PATCH("/:id/toggle", controllers.ToggleTaxStatus)
	}
}
