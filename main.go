package main

// @title           B-Resto API
// @version         1.0
// @description     API RESTful para sistema POS de restaurantes
// @termsOfService  http://swagger.io/terms/

// @contact.name   Soporte API
// @contact.email  support@b-resto.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Escribir "Bearer {token}" para autenticación JWT

import (
	"b-resto/config"
	"b-resto/models"
	"b-resto/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err := db.AutoMigrate(
		// Existentes
		&models.User{},
		&models.Unit{},
		&models.Company{},
		&models.PaymentMethod{},
		&models.Tax{},

		// Nuevos - Base
		&models.InventoryCategory{},
		&models.ProductCategory{},
		&models.KitchenStation{},
		&models.Warehouse{},
		&models.TableArea{},
		&models.Table{},

		// Nuevos - Productos
		&models.ProductTemplate{},
		&models.ProductProduct{},
		&models.ProductAttribute{},
		&models.ProductAttributeValue{},

		// Nuevos - Combos
		&models.Combo{},
		&models.ComboItem{},

		// Nuevos - Órdenes y Ventas
		&models.Sequence{},
		&models.Journal{},
		&models.Order{},
		&models.OrderItem{},
		&models.OrderPayment{},
		&models.KitchenTicket{},
		&models.KitchenTicketItem{},

		// Nuevos - POS y Caja
		&models.POS{},
		&models.POSSession{},
		&models.CashMovement{},

		// Nuevos - Inventario
		&models.Partner{},
		&models.Recipe{},
		&models.StockTransfer{},
		&models.StockTransferItem{},
		&models.PurchaseOrder{},
		&models.PurchaseOrderItem{},
		&models.Inventory{},

		// Nuevos - Reservaciones
		&models.Reservation{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("✅ Database migration completed successfully")

	config.DB = db

	config.InitCasbin()
	config.SeedCasbinPolicies()

	r := gin.Default()
	r.Use(CORSMiddleware())
	routes.SetupRoutes(r)

	log.Printf("Server starting on port %s", config.ServerPort)
	log.Printf("Environment: %s", gin.Mode())

	if err := r.Run(config.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
