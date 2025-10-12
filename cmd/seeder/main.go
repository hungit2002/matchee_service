package main

import (
	"log"

	"matchee/services/internal/config"
	"matchee/services/internal/entity"
	"matchee/services/pkg/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := database.NewMySQL(cfg)
	if err != nil {
		log.Fatalf("mysql: %v", err)
	}

	// Auto migrate entities
	if err := db.AutoMigrate(&entity.Role{}); err != nil {
		log.Fatalf("auto migrate: %v", err)
	}

	// Seed roles
	roles := []entity.Role{
		{Name: "admin"},
		{Name: "user"},
		{Name: "moderator"},
		{Name: "premium"},
	}

	for _, role := range roles {
		var existingRole entity.Role
		if err := db.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			// Role doesn't exist, create it
			if err := db.Create(&role).Error; err != nil {
				log.Printf("failed to create role %s: %v", role.Name, err)
			} else {
				log.Printf("created role: %s", role.Name)
			}
		} else {
			log.Printf("role %s already exists", role.Name)
		}
	}

	log.Println("seeder completed")
}
