package routes

import (
	"b-resto/controllers"

	"github.com/gin-gonic/gin"
)

// SetupTableAreaRoutes configura las rutas para table areas
func SetupTableAreaRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/table-areas", controllers.GetTableAreas)
		api.GET("/table-areas/:id", controllers.GetTableArea)
		api.POST("/table-areas", controllers.CreateTableArea)
		api.PUT("/table-areas/:id", controllers.UpdateTableArea)
		api.DELETE("/table-areas/:id", controllers.DeleteTableArea)
		api.PATCH("/table-areas/:id/toggle", controllers.ToggleTableAreaStatus)
	}
}
