package emailalert

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

var srv *gmail.Service
var config *oauth2.Config
var Enabled bool
var to string = "dmitriy8096@gmail.com"

func EmailAlertStart(Enabled_ bool, to_ string) {
	Enabled = Enabled_
	if err, res := CheckConfig(); res && Enabled {
		to = to_
		GetClientToken()
	} else {
		log.Println(err)
		Enabled = false
	}

}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "emailalert/token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(oauth2.NoContext, tok)

}

// Request a token from the web, then returns the retrieved token.
////////      Rewrite for the frontend /////////
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {

	authURL := config.AuthCodeURL("CSRF", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func CheckConfig() (msgerror string, status bool) {
	var err error
	b, err := ioutil.ReadFile("emailalert/credentials.json")
	if err != nil {
		return "Unable to read client secret file: " + err.Error(), false
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err = google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		config = nil
		return "Unable to parse client secret file to config: " + err.Error(), false
	}
	return "Ok", true
}

func GetClientToken() {
	client := getClient(config)
	var err error
	srv, err = gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
}

func SendEmailAlert(subject, msg string) {
	if Enabled {
		var message gmail.Message
		var err error
		// Compose the message
		messageStr := []byte(

			"To: " + to + "\r\n" +
				"Subject: " + subject + "\r\n\r\n" +
				msg)

		// Place messageStr into message.Raw in base64 encoded format
		message.Raw = base64.URLEncoding.EncodeToString(messageStr)

		// Send the message
		_, err = srv.Users.Messages.Send("me", &message).Do()
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Println("Message sent!")
		}
	}
}
