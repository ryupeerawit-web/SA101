package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	DBDSN string
}

func LoadConfig() *Config {
	// พยายามโหลด .env จาก root directory หรือ parent directory
	err := godotenv.Load()
	if err != nil {
		// ลองอ่านจาก path สำรองเผื่อรันจาก cmd/sever
		_ = godotenv.Load("../../.env")
		log.Println("Notice: No primary .env file found, using fallback or system env")
	}

	return &Config{
		Port:  getEnv("PORT", "8080"),
		DBDSN: getEnv("DATABASE_URL", "host=localhost user=postgres password=postgrespassword dbname=fitness_db port=5432 sslmode=disable"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}