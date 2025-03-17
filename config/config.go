package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

const (
	credentialsFile = "credentials.json"
	tokenFileName   = "token.json"
	authCodeState   = "state-token"
)

var oauthConfig *oauth2.Config

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetOAuthConfig() *oauth2.Config {
	if oauthConfig == nil {
		loadEnv()

		oauthConfig = &oauth2.Config{
			ClientID:     os.Getenv("GMAIL_CLIENT_ID"),
			ClientSecret: os.Getenv("GMAIL_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GMAIL_REDIRECT_URI"),
			Scopes: []string{
				gmail.GmailReadonlyScope,
				gmail.GmailModifyScope,
				gmail.GmailLabelsScope,
			},
			Endpoint: google.Endpoint,
		}

		if oauthConfig.ClientID == "" || oauthConfig.ClientSecret == "" || oauthConfig.RedirectURL == "" {
			log.Fatal("GMAIL API credentials are missing in .env file")
		}
	}

	return oauthConfig
}

func loadToken() (*oauth2.Token, error) {
	file, err := os.Open(tokenFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	tok := &oauth2.Token{}
	err = json.NewDecoder(file).Decode(tok)

	return tok, err
}

func saveToken(token *oauth2.Token) {
	file, err := os.Create(tokenFileName)
	if err != nil {
		log.Fatalf("Failed to save token: %v", err)
	}

	defer file.Close()

	json.NewEncoder(file).Encode(token)
	fmt.Println("Token successfully saved to", tokenFileName)
}

func getTokenFromWeb() *oauth2.Token {
	config := GetOAuthConfig()
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Println("Go to the following link to authorize:", authURL)

	var authCode string

	fmt.Print("Enter the code: ")
	fmt.Scan(&authCode)

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Failed to exchange auth code for token: %v", err)
	}

	return tok
}

func GetGmailClient() *oauth2.Token {
	tok, err := loadToken()
	if err == nil {
		return tok
	}

	tok = getTokenFromWeb()
	saveToken(tok)

	return tok
}

func GetGmailService() (*gmail.Service, error) {
	config := GetOAuthConfig()
	tok := GetGmailClient()

	ctx := context.Background()
	tokenSource := config.TokenSource(ctx, tok)
	oauthClient := oauth2.NewClient(ctx, tokenSource)

	refreshedToken, err := tokenSource.Token()
	if err == nil {
		saveToken(refreshedToken)
	}

	return gmail.NewService(context.Background(), option.WithHTTPClient(oauthClient))
}
