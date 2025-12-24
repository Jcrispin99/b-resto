package routes

import (
	"b-resto/config"
	"b-resto/controllers"
	"b-resto/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "b-resto/docs" // SetupRoutes configura todas las rutas de la aplicación
)

func SetupRoutes(r *gin.Engine) {
	// Ruta raíz
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to B-Resto API"})
	})

	// Swagger UI (solo en desarrollo)
	if config.GetEnvironment() == "development" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Rutas de autenticación
	auth := r.Group("/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
	}

	// Rutas públicas
	public := r.Group("/public")
	{
		public.GET("/data", controllers.GetPublicData)
	}

	// Rutas protegidas con JWT + Casbin
	api := r.Group("/api")

	// Solo aplicar middlewares de seguridad en modo producción
	if config.GetEnvironment() != "development" {
		api.Use(middlewares.AuthMiddleware())
		api.Use(middlewares.CasbinMiddleware())
	}

	{
		api.GET("/profile", controllers.GetProfile)

		api.GET("/resources", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "List of resources",
				"resources": []gin.H{
					{"id": 1, "name": "Resource 1"},
					{"id": 2, "name": "Resource 2"},
					{"id": 3, "name": "Resource 3"},
				},
			})
		})

		api.PUT("/resources/:id", controllers.UpdateResource)

		// Admin routes (protegidas por políticas de Casbin)
		admin := api.Group("/admin")
		{
			admin.POST("/resources", controllers.CreateResource)
			admin.DELETE("/resources/:id", controllers.DeleteResource)
			admin.GET("/users", controllers.GetUsers)
		}

		// Rutas API - endpoints
		SetupUnitsRoutes(api)
		SetupCompanyRoutes(api)
		SetupTaxesRoutes(api)
		SetupPaymentMethodsRoutes(api)

		// FASE 1: Catálogos Base
		SetupWarehouseRoutes(r)
		SetupKitchenStationRoutes(r)
		SetupInventoryCategoryRoutes(r)

		// FASE 2: Categorías de Productos
		SetupProductCategoryRoutes(r)

		// FASE 5: Secuencias y Journals
		SetupSequenceRoutes(r)
		SetupJournalRoutes(r)

		// FASE 6: Mesas y Reservaciones
		SetupTableAreaRoutes(r)
		SetupTableRoutes(r)
		SetupReservationRoutes(r)

		// FASE 7: Órdenes de Venta (CRÍTICO POS)
		SetupOrderRoutes(r)
		SetupKitchenTicketRoutes(r)

		// FASE 8: POS y Caja
		SetupPOSRoutes(r)
		SetupPOSSessionRoutes(r)

		// FASE 9: Proveedores/Clientes
		SetupPartnerRoutes(r)

		// FASE 10: Compras y Transferencias
		SetupPurchaseOrderRoutes(r)
		SetupStockTransferRoutes(r)

		// FASE 11: Inventario (Kardex)
		SetupInventoryRoutes(r)

		// FASE 3-4: Módulo de Productos (COMPLETO)
		SetupProductRoutes(r)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "auth-go-gin",
		})
	})
}
