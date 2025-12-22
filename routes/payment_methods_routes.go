package routes

import (
	"b-resto/controllers"

	"github.com/gin-gonic/gin"
)

func SetupPaymentMethodsRoutes(router *gin.RouterGroup) {
	paymentMethods := router.Group("/payment-methods")
	{
		paymentMethods.GET("", controllers.GetPaymentMethods)
		paymentMethods.GET("/:id", controllers.GetPaymentMethod)
		paymentMethods.POST("", controllers.CreatePaymentMethod)
		paymentMethods.PUT("/:id", controllers.UpdatePaymentMethod)
		paymentMethods.DELETE("/:id", controllers.DeletePaymentMethod)
		paymentMethods.PATCH("/:id/toggle", controllers.TogglePaymentMethodStatus)
	}
}
