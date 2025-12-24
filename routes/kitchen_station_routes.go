package routes

import (
	"b-resto/controllers"

	"github.com/gin-gonic/gin"
)

// SetupKitchenStationRoutes configura las rutas para kitchen stations
func SetupKitchenStationRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/kitchen-stations", controllers.GetKitchenStations)
		api.GET("/kitchen-stations/:id", controllers.GetKitchenStation)
		api.POST("/kitchen-stations", controllers.CreateKitchenStation)
		api.PUT("/kitchen-stations/:id", controllers.UpdateKitchenStation)
		api.DELETE("/kitchen-stations/:id", controllers.DeleteKitchenStation)
		api.PATCH("/kitchen-stations/:id/toggle", controllers.ToggleKitchenStationStatus)
	}
}
