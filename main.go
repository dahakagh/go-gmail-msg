package main

import (
	"go-gmail-msg/gmail"
	"go-gmail-msg/services"
	"log"
)

func main() {
	srv, err := gmail.GetGmailService()
	if err != nil {
		log.Fatalf("Failed to initialize Gmail API: %v", err)
	}

	services.FetchUnreadEmails(srv)
}
