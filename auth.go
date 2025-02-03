package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
)

func exchangeCodeForToken(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", config.ClientID)
	data.Set("client_secret", config.ClientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", config.RedirectURI)

	resp, err := http.PostForm(config.TokenEndpoint, data)
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", errors.New("Error reading response body: " + err.Error())
		}
		msg := "Error making POST request: " + string(body)
		log.Print("ERROR: ", msg)
		return "", errors.New(msg)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	log.Println("Token response:", string(body))
	return string(body), nil
}

func redirectToAuthenticate(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", config.AuthenticationEndpoint, nil)
	if err != nil {
		log.Fatal("Error creating authorization request: ", err)
	}
	query := req.URL.Query()
	query.Add("client_id", config.ClientID)
	query.Add("response_type", "code")
	query.Add("redirect_uri", config.RedirectURI)
	query.Add("response_mode", "query")
	query.Add("scope", config.Scope)
	// TODO: state and nonce from config/generated
	query.Add("state", "12345")
	query.Add("nonce", "678910")

	req.URL.RawQuery = query.Encode()

	// Redirect the user to the Azure AD login page
	http.Redirect(w, r, req.URL.String(), http.StatusFound)
}
