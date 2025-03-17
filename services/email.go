package services

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"google.golang.org/api/gmail/v1"
)

const outputDir = "emails"

func FetchUnreadEmails(srv *gmail.Service) {
	user := "me"
	query := "is:unread"

	r, err := srv.Users.Messages.List(user).Q(query).Do()
	if err != nil {
		log.Fatalf("Failed to fetch emails: %v", err)
	}

	if len(r.Messages) == 0 {
		fmt.Println("No unread emails found.")
		return
	}

	os.MkdirAll(outputDir, os.ModePerm)

	for _, msg := range r.Messages {
		message, err := srv.Users.Messages.Get(user, msg.Id).Format("full").Do()
		if err != nil {
			log.Printf("Failed to retrieve email %s: %v", msg.Id, err)
			continue
		}

		var subject string

		for _, header := range message.Payload.Headers {
			if header.Name == "Subject" {
				subject = header.Value
				break
			}
		}

		body := extractMessageBody(message.Payload.Parts)

		filename := fmt.Sprintf("%s/%s.txt", outputDir, msg.Id)

		err = os.WriteFile(filename, []byte(subject+"\n\n"+body), 0644)
		if err != nil {
			log.Printf("Failed to save email %s: %v", msg.Id, err)
		} else {
			fmt.Printf("Email saved: %s\n", filename)
		}

		SaveAttachments(srv, message)
		MarkAsRead(srv, msg.Id)
	}
}

func extractMessageBody(parts []*gmail.MessagePart) string {
	var body string

	for _, part := range parts {
		if part.MimeType == "text/plain" || part.MimeType == "text/html" {
			data, err := base64.URLEncoding.DecodeString(part.Body.Data)
			if err == nil {
				body += string(data) + "\n"
			}
		} else if len(part.Parts) > 0 {
			body += extractMessageBody(part.Parts)
		}
	}
	return body
}

func MarkAsRead(srv *gmail.Service, msgID string) {
	_, err := srv.Users.Messages.Modify("me", msgID, &gmail.ModifyMessageRequest{
		RemoveLabelIds: []string{"UNREAD"},
	}).Do()
	if err != nil {
		log.Printf("Failed to mark email %s as read: %v", msgID, err)
	} else {
		fmt.Printf("Email %s marked as read.\n", msgID)
	}
}
