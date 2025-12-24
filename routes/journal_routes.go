package routes

import (
	"b-resto/controllers"

	"github.com/gin-gonic/gin"
)

// SetupJournalRoutes configura las rutas para journals
func SetupJournalRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/journals", controllers.GetJournals)
		api.GET("/journals/:id", controllers.GetJournal)
		api.POST("/journals", controllers.CreateJournal)
		api.PUT("/journals/:id", controllers.UpdateJournal)
		api.PATCH("/journals/:id/toggle", controllers.ToggleJournalStatus)
	}
}
