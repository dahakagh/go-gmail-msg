package services

import (
	"encoding/base64"
	"fmt"
	"go-gmail-msg/utils"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"google.golang.org/api/gmail/v1"
)

const outputDir = "emails"

func FetchUnreadEmails(service *gmail.Service) {
	user := "me"
	query := "is:unread"

	response, err := service.Users.Messages.List(user).Q(query).Do()
	if err != nil {
		log.Fatalf("Failed to fetch emails: %v", err)
	}

	if len(response.Messages) == 0 {
		fmt.Println("No unread emails found.")
		return
	}

	os.MkdirAll(outputDir, os.ModePerm)

	for _, msg := range response.Messages {
		message, err := service.Users.Messages.Get(user, msg.Id).Format("full").Do()
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

		date := extractEmailDate(message)
		from := extractSenderEmail(message)

		fileName := fmt.Sprintf("%s_%s.txt", date, from)
		filePath := fmt.Sprintf("%s/%s", outputDir, fileName)

		body := extractMessageBody(message.Payload.Parts)

		err = os.WriteFile(filePath, []byte(subject+"\n\n"+body), 0644)
		if err != nil {
			log.Printf("Failed to save email %s: %v", msg.Id, err)
		} else {
			fmt.Printf("Email saved: %s\n", fileName)
		}

		SaveAttachments(service, message)
		MarkAsRead(service, msg.Id)
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

func extractEmailDate(msg *gmail.Message) string {
	if msg == nil || msg.Payload == nil {
		return "unknown-date"
	}

	for _, header := range msg.Payload.Headers {
		if strings.EqualFold(header.Name, "Date") {
			dateStr := header.Value

			dateStr = strings.TrimSuffix(dateStr, " (UTC)")

			formats := []string{
				time.RFC1123Z,
				time.RFC1123,
			}

			for _, format := range formats {
				parsedTime, err := time.Parse(format, dateStr)
				if err == nil {
					return parsedTime.Format("2006-01-02")
				}
			}

			log.Printf("Failed to parse date after cleanup: %s", dateStr)
			return "unknown-date"
		}
	}

	return "unknown-date"
}

func extractSenderEmail(msg *gmail.Message) string {
	if msg == nil || msg.Payload == nil {
		return "unknown"
	}

	for _, header := range msg.Payload.Headers {
		if header.Name == "From" {
			re := regexp.MustCompile(`<([^>]+)>`)
			matches := re.FindStringSubmatch(header.Value)

			if len(matches) > 1 {
				return utils.SanitizeFileName(matches[1])
			}

			return utils.SanitizeFileName(header.Value)
		}
	}
	return "unknown"
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
