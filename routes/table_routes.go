package routes

import (
	"b-resto/controllers"

	"github.com/gin-gonic/gin"
)

// SetupTableRoutes configura las rutas para tables
func SetupTableRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/tables", controllers.GetTables)
		api.GET("/tables/:id", controllers.GetTable)
		api.POST("/tables", controllers.CreateTable)
		api.PUT("/tables/:id", controllers.UpdateTable)
		api.DELETE("/tables/:id", controllers.DeleteTable)
		api.PATCH("/tables/:id/toggle", controllers.ToggleTableStatus)
		api.GET("/tables/area/:area_id", controllers.GetTablesByArea)
		api.GET("/tables/available", controllers.GetAvailableTables)
	}
}
