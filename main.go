package main

import (
	"go-gmail-msg/config"
	"go-gmail-msg/services"
	"log"
)

func main() {
	srv, err := config.GetGmailService()
	if err != nil {
		log.Fatalf("Failed to initialize Gmail API: %v", err)
	}

	services.FetchUnreadEmails(srv)
}
