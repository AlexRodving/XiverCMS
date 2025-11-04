package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	GinMode        string
	DBDriver       string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBPath         string
	JWTSecret      string
	JWTExpiration  string
	CORSOrigin     string
	AllowedOrigins []string
}

var AppConfig *Config

func LoadConfig() {
	// Load .env file if it exists
	_ = godotenv.Load()

	AppConfig = &Config{
		Port:           getEnv("PORT", "8080"),
		GinMode:        getEnv("GIN_MODE", "debug"),
		DBDriver:       getEnv("DB_DRIVER", "sqlite"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "xivercrm"),
		DBPath:         getEnv("DB_PATH", "./data/xivercrm.db"),
		JWTSecret:      getEnv("JWT_SECRET", "change-this-secret-key-in-production"),
		JWTExpiration:  getEnv("JWT_EXPIRATION", "24h"),
		CORSOrigin:     getEnv("CORS_ORIGIN", "http://localhost:5173"),
		AllowedOrigins: getEnvArray("ALLOWED_ORIGINS", []string{"http://localhost:5173", "http://localhost:3000"}),
	}

	log.Println("Configuration loaded successfully")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvArray(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	// Simple comma-separated parsing
	result := []string{}
	start := 0
	for i, char := range value {
		if char == ',' {
			if i > start {
				result = append(result, value[start:i])
			}
			start = i + 1
		}
	}
	if start < len(value) {
		result = append(result, value[start:])
	}

	if len(result) == 0 {
		return defaultValue
	}
	return result
}

func GetPort() string {
	return AppConfig.Port
}

func GetDBDriver() string {
	return AppConfig.DBDriver
}

func GetDBConnectionString() string {
	if AppConfig.DBDriver == "sqlite" {
		return AppConfig.DBPath
	}

	port, _ := strconv.Atoi(AppConfig.DBPort)
	if port == 0 {
		port = 5432
	}

	return "host=" + AppConfig.DBHost +
		" port=" + strconv.Itoa(port) +
		" user=" + AppConfig.DBUser +
		" password=" + AppConfig.DBPassword +
		" dbname=" + AppConfig.DBName +
		" sslmode=disable"
}
