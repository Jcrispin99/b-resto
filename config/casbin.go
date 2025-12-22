package config

import (
	"log"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

var Enforcer *casbin.Enforcer

func InitCasbin() {
	// Usar el adaptador GORM (guarda políticas en PostgreSQL)
	adapter, err := gormadapter.NewAdapterByDB(DB)
	if err != nil {
		log.Fatalf("Failed to create Casbin adapter: %v", err)
	}

	// Cargar el modelo
	enforcer, err := casbin.NewEnforcer("config/casbin_model.conf", adapter)
	if err != nil {
		log.Fatalf("Failed to create Casbin enforcer: %v", err)
	}

	// Cargar políticas desde la BD
	enforcer.LoadPolicy()

	Enforcer = enforcer
	log.Println("✅ Casbin initialized successfully")
}

// SeedCasbinPolicies seedea políticas iniciales
func SeedCasbinPolicies() {
	// Asignar roles a usuarios
	Enforcer.AddRoleForUser("admin", "admin_role")
	Enforcer.AddRoleForUser("user", "user_role")

	// Permisos para admin_role - Acceso total a /api
	Enforcer.AddPolicy("admin_role", "/api/*", "*")

	// Permisos para user_role
	Enforcer.AddPolicy("user_role", "/api/profile", "GET")
	Enforcer.AddPolicy("user_role", "/api/units", "GET")
	Enforcer.AddPolicy("user_role", "/api/units/*", "GET")

	// Guardar cambios
	Enforcer.SavePolicy()
	log.Println("✅ Casbin policies seeded")
}
