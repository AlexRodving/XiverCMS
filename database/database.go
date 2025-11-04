package database

import (
	"log"
	"os"

	"github.com/xivercms/xivercms/config"
	"github.com/xivercms/xivercms/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	var err error
	var dialector gorm.Dialector

	dbDriver := config.GetDBDriver()

	switch dbDriver {
	case "postgres":
		dsn := config.GetDBConnectionString()
		dialector = postgres.Open(dsn)
	case "sqlite":
		// Create data directory if it doesn't exist
		dbPath := config.AppConfig.DBPath
		dir := "./data"
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0755)
		}
		dialector = sqlite.Open(dbPath)
	default:
		log.Fatal("Unsupported database driver: ", dbDriver)
	}

	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")
}

func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.ContentType{},
		&models.ContentEntry{},
		&models.MediaFile{},
		&models.APIToken{},
		&models.AuditLog{},
		&models.ContentHistory{},
		&models.ContentRelation{},
		&models.ComponentType{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed")
}

func Seed() {
	// Check if admin user already exists
	var adminUser models.User
	if err := DB.Where("email = ?", "admin@xivercms.com").First(&adminUser).Error; err == nil {
		log.Println("Admin user already exists, skipping seed")
		return
	}

	// Create default admin user
	hashedPassword, _ := hashPassword("admin123") // In production, use a secure password
	adminUser = models.User{
		Email:        "admin@xivercms.com",
		Username:     "admin",
		Password:     hashedPassword,
		FirstName:    "Admin",
		LastName:     "User",
		IsActive:     true,
		IsSuperAdmin: true,
	}

	if err := DB.Create(&adminUser).Error; err != nil {
		log.Printf("Failed to create admin user: %v", err)
		return
	}

	// Create default roles
	roles := []models.Role{
		{
			Name:        "Public",
			Description: "Public role",
			Type:        "public",
		},
		{
			Name:        "Authenticated",
			Description: "Authenticated user role",
			Type:        "public",
		},
		{
			Name:        "Admin",
			Description: "Administrator role",
			Type:        "custom",
		},
	}

	for _, role := range roles {
		var existingRole models.Role
		if err := DB.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			DB.Create(&role)
		}
	}

	log.Println("Database seeded successfully")
}

func hashPassword(password string) (string, error) {
	// This is a placeholder - you should use bcrypt or similar
	// For now, we'll implement it in the auth package
	return password, nil
}
