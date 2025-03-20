package gmail

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"go-gmail-msg/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

const (
	tokenFileName = "token.json"
	authCodeState = "state-token"
)

var oauthConfig *oauth2.Config

func GetOAuthConfig() *oauth2.Config {
	if oauthConfig == nil {
		config.LoadEnv()

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

		if oauthConfig.ClientID == "" || oauthConfig.ClientSecret == "" || oauthConfig.RedirectURL == "" || os.Getenv("HTTP_SERVER_ADDRESS") == "" {
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

	token := &oauth2.Token{}
	err = json.NewDecoder(file).Decode(token)

	return token, err
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
	host := os.Getenv("HTTP_SERVER_ADDRESS")

	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal("Server startup error:", err)
	}

	defer listener.Close()

	config := GetOAuthConfig()
	authURL := config.AuthCodeURL(authCodeState, oauth2.AccessTypeOffline)

	fmt.Println("Go to the following link to authorize:", authURL)

	authCodeChannel := make(chan string)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		authCodeChannel <- r.URL.Query().Get("code")

		fmt.Fprintf(w, "Authentication successful, you can close the window.")
	})

	go http.Serve(listener, nil)

	authCode := <-authCodeChannel

	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Failed to exchange auth code for token: %v", err)
	}

	return token
}

func GetGmailClient() *oauth2.Token {
	token, err := loadToken()
	if err == nil {
		return token
	}

	token = getTokenFromWeb()
	saveToken(token)

	return token
}

func GetGmailService() (*gmail.Service, error) {
	config := GetOAuthConfig()
	tok := GetGmailClient()

	ctx := context.Background()
	tokenSource := config.TokenSource(ctx, tok)
	oauthClient := oauth2.NewClient(ctx, tokenSource)

	refreshedToken, err := tokenSource.Token()
	if err == nil && refreshedToken.AccessToken != tok.AccessToken {
		saveToken(refreshedToken)
	}

	return gmail.NewService(context.Background(), option.WithHTTPClient(oauthClient))
}
