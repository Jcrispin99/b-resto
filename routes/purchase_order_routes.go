package routes

import (
	"b-resto/controllers"

	"github.com/gin-gonic/gin"
)

// SetupPurchaseOrderRoutes configura las rutas para purchase orders
func SetupPurchaseOrderRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/purchase-orders", controllers.GetPurchaseOrders)
		api.GET("/purchase-orders/:id", controllers.GetPurchaseOrder)
		api.POST("/purchase-orders", controllers.CreatePurchaseOrder)
		api.PUT("/purchase-orders/:id", controllers.UpdatePurchaseOrder)
		api.PATCH("/purchase-orders/:id/send", controllers.SendPurchaseOrder)
		api.PATCH("/purchase-orders/:id/receive", controllers.ReceivePurchaseOrder)
		api.PATCH("/purchase-orders/:id/cancel", controllers.CancelPurchaseOrder)
	}
}
