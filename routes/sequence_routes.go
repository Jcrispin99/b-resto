package routes

import (
	"b-resto/controllers"

	"github.com/gin-gonic/gin"
)

// SetupSequenceRoutes configura las rutas para sequences
func SetupSequenceRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/sequences", controllers.GetSequences)
		api.GET("/sequences/:id", controllers.GetSequence)
		api.POST("/sequences", controllers.CreateSequence)
		api.PUT("/sequences/:id", controllers.UpdateSequence)
		api.GET("/sequences/:id/next", controllers.GetNextNumber)
	}
}
