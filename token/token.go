package token

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/kelr/gundyr/auth"
	"golang.org/x/oauth2"
)

const tokenFile = "token.json"

// Logic .
type Logic struct {
	clientID     string
	clientSecret string
}

// New .
func New(clientID, clientSecret string) *Logic {
	return &Logic{
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

// Get .
func (l *Logic) Get() (*oauth2.Token, error) {
	scopes := []string{"chat:read", "chat:edit"}

	// Setup OAuth2 configuration
	config, err := auth.NewUserAuth(l.clientID, l.clientSecret, "http://localhost:8080", &scopes)
	if err != nil {
		return nil, err
	}

	var token *oauth2.Token
	if _, err := os.Stat(tokenFile); errors.Is(err, os.ErrNotExist) {
		token, err = generateNewToken(config)
		if err != nil {
			return nil, err
		}

		// Write the token to a file
		if err := auth.FlushTokenFile(tokenFile, token); err != nil {
			return nil, err
		}
	} else {
		token, err = auth.RetrieveTokenFile(config, tokenFile)
		if err != nil {
			return nil, err
		}
	}

	return token, nil
}

func startCallbackServer(authCodeCh chan<- string) func() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "Code received, you can close this window now.")
		authCodeCh <- r.URL.Query().Get("code")
	})

	server := &http.Server{Addr: "localhost:8080", Handler: mux}
	go server.ListenAndServe()

	return func() {
		_ = server.Close()
	}
}

func generateNewToken(config *oauth2.Config) (*oauth2.Token, error) {
	// Get the URL to send to the user and the state code to protect against CSRF attacks.
	url, state := auth.GetAuthCodeURL(config)
	fmt.Println("State generated:", state)
	fmt.Println("Open the URL and do the auth please:")
	fmt.Println(url)
	// fmt.Println("Ensure that state received at URI is:", state)
	// fmt.Println("Copy the authCode from the browser")

	// Enter the code received by the redirect URI. Ensure that the state value
	// obtained at the redirect URI matches the previous state value.
	// var authCode string
	// if _, err := fmt.Scan(&authCode); err != nil {
	// 	return nil, err
	// }

	// TODO: Also read and compare the state!
	authCodeCh := make(chan string)
	close := startCallbackServer(authCodeCh)
	authCode := <-authCodeCh
	close()

	// Obtain the user token through the code. This token can be reused as long as
	// it has not expired, but the auth code cannot be reused.
	token, err := auth.TokenExchange(config, authCode)
	if err != nil {
		return nil, err
	}
	return token, nil
}
