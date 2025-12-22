package routes

import (
	"b-resto/controllers"

	"github.com/gin-gonic/gin"
)

// SetupUnitsRoutes configura las rutas de units
func SetupUnitsRoutes(router *gin.RouterGroup) {
	units := router.Group("/units")
	{
		units.GET("", controllers.GetUnits)
		units.GET("/:id", controllers.GetUnit)
		units.POST("", controllers.CreateUnit)
		units.PUT("/:id", controllers.UpdateUnit)
		units.DELETE("/:id", controllers.DeleteUnit)
		units.PATCH("/:id/toggle", controllers.ToggleUnitStatus)
	}
}
