package routes

import (
	"b-resto/controllers"

	"github.com/gin-gonic/gin"
)

// SetupReservationRoutes configura las rutas para reservations
func SetupReservationRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/reservations", controllers.GetReservations)
		api.GET("/reservations/:id", controllers.GetReservation)
		api.POST("/reservations", controllers.CreateReservation)
		api.PUT("/reservations/:id", controllers.UpdateReservation)
		api.PATCH("/reservations/:id/confirm", controllers.ConfirmReservation)
		api.PATCH("/reservations/:id/cancel", controllers.CancelReservation)
		api.PATCH("/reservations/:id/seat", controllers.SeatReservation)
		api.GET("/reservations/date/:date", controllers.GetReservationsByDate)
	}
}
