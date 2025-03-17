package services

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"google.golang.org/api/gmail/v1"
)

const attachmentsDir = "attachments"

func SaveAttachments(srv *gmail.Service, msg *gmail.Message) {
	os.MkdirAll(attachmentsDir, os.ModePerm)

	for _, part := range msg.Payload.Parts {
		if part.Filename != "" {
			attID := part.Body.AttachmentId

			att, err := srv.Users.Messages.Attachments.Get("me", msg.Id, attID).Do()
			if err != nil {
				log.Printf("Failed to fetch attachment %s: %v", part.Filename, err)
				continue
			}

			data, err := base64.URLEncoding.DecodeString(att.Data)
			if err != nil {
				log.Printf("Failed to decode attachment %s: %v", part.Filename, err)
				continue
			}

			filePath := fmt.Sprintf("%s/%s", attachmentsDir, part.Filename)

			err = os.WriteFile(filePath, data, 0644)
			if err != nil {
				log.Printf("Failed to save attachment %s: %v", part.Filename, err)
			} else {
				fmt.Printf("Attachment saved: %s\n", filePath)
			}
		}
	}
}
