package config

import (
	"fmt"
	"log"
	"os"

	"go-gmail-msg/tui"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	isEnvFileExists := true
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		fmt.Println("⚠️  .env file not found. Creating a new one...")
		isEnvFileExists = false
	}

	_ = godotenv.Load()

	requiredKeys := []string{
		"GMAIL_CLIENT_ID",
		"GMAIL_CLIENT_SECRET",
		"GMAIL_REDIRECT_URI",
		"HTTP_SERVER_ADDRESS",
	}

	envData := make(map[string]string)
	for _, key := range requiredKeys {
		if os.Getenv(key) == "" {
			envData[key] = tui.PromptUser(key)
		}
	}

	if !isEnvFileExists || len(envData) > 0 {
		saveEnvFile(envData)

		if err := godotenv.Load(); err != nil {
			log.Fatal("Error reloading .env file:", err)
		}
	}
}

func saveEnvFile(envData map[string]string) {
	file, err := os.OpenFile(".env", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error creating .env file: %v", err)
	}
	defer file.Close()

	for key, value := range envData {
		_, err := file.WriteString(fmt.Sprintf("%s=%s\n", key, value))
		if err != nil {
			log.Fatalf("Error writing to .env file: %v", err)
		}
	}

	fmt.Println("✅ .env file has been successfully updated.")
}
