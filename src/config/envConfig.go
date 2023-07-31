package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig(fileName ...string) {
	var envName string
	if len(fileName) == 0 {
		envName = ".env"
	} else {
		envName = fileName[0]

	}

	err := godotenv.Load(envName)

	if err != nil {
		log.Println(err)
		log.Fatal("Error while loading env variables.")
	}
}

func GetConfig(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
