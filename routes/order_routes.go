package routes

import (
	"b-resto/controllers"

	"github.com/gin-gonic/gin"
)

// SetupOrderRoutes configura las rutas para orders
func SetupOrderRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/orders", controllers.GetOrders)
		api.POST("/orders", controllers.CreateOrder)

		// Rutas específicas ANTES de las dinámicas con :id
		api.GET("/orders/table/:table_id", controllers.GetOrdersByTable)

		// Rutas dinámicas con :id
		api.GET("/orders/:id", controllers.GetOrder)
		api.PUT("/orders/:id", controllers.UpdateOrder)
		api.PATCH("/orders/:id/confirm", controllers.ConfirmOrder)
		api.PATCH("/orders/:id/cancel", controllers.CancelOrder)
		api.PATCH("/orders/:id/complete", controllers.CompleteOrder)

		// Pagos de orden - usar :id consistentemente
		api.GET("/orders/:id/payments", controllers.GetOrderPayments)
		api.POST("/orders/:id/payments", controllers.CreateOrderPayment)
		api.DELETE("/orders/:id/payments/:payment_id", controllers.DeleteOrderPayment)
	}
}
