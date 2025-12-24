package routes

import (
	"b-resto/controllers"

	"github.com/gin-gonic/gin"
)

// SetupPartnerRoutes configura las rutas para partners
func SetupPartnerRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/partners", controllers.GetPartners)
		api.GET("/partners/:id", controllers.GetPartner)
		api.POST("/partners", controllers.CreatePartner)
		api.PUT("/partners/:id", controllers.UpdatePartner)
		api.DELETE("/partners/:id", controllers.DeletePartner)
		api.PATCH("/partners/:id/toggle", controllers.TogglePartnerStatus)
		api.GET("/partners/suppliers", controllers.GetSuppliers)
		api.GET("/partners/customers", controllers.GetCustomers)
	}
}
