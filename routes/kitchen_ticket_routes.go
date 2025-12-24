package routes

import (
	"b-resto/controllers"

	"github.com/gin-gonic/gin"
)

// SetupKitchenTicketRoutes configura las rutas para kitchen tickets
func SetupKitchenTicketRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/kitchen-tickets", controllers.GetKitchenTickets)
		api.GET("/kitchen-tickets/:id", controllers.GetKitchenTicket)
		api.PATCH("/kitchen-tickets/:id/status", controllers.UpdateTicketStatus)
		api.PATCH("/kitchen-tickets/:id/preparing", controllers.MarkTicketPreparing)
		api.PATCH("/kitchen-tickets/:id/ready", controllers.MarkTicketReady)
		api.PATCH("/kitchen-tickets/:id/delivered", controllers.MarkTicketDelivered)
		api.GET("/kitchen-tickets/station/:station_id", controllers.GetTicketsByStation)
	}
}
