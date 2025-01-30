package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func exchangeCodeForToken(code string) (string, error) {
	tokenEndpoint := "https://login.microsoftonline.com/3aa4a235-b6e2-48d5-9195-7fcf05b459b0/oauth2/v2.0/token"
	clientID := "5cb6c4de-28d0-4b62-a547-262dc2377baf"
	clientSecret := "" // TODO: Get from envvars
	redirectURI := "http://localhost:8080/callback"

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	resp, err := http.PostForm(tokenEndpoint, data)
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", errors.New("Error reading response body: " + err.Error())
		}
		msg := "Error making POST request: " + string(body)
		log.Print("ERROR: ", msg)
		return "", errors.New(msg)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	log.Println("Token response:", string(body))
	return string(body), nil
}

func redirectToAuthenticate(w http.ResponseWriter, r *http.Request) {
	clientID := "5cb6c4de-28d0-4b62-a547-262dc2377baf"
	redirectURI := "http://localhost:8080/callback"
	authorizeEndpoint := "https://login.microsoftonline.com/3aa4a235-b6e2-48d5-9195-7fcf05b459b0/oauth2/v2.0/authorize"

	req, err := http.NewRequest("GET", authorizeEndpoint, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	query := req.URL.Query()
	query.Add("client_id", clientID)
	query.Add("response_type", "code")
	query.Add("redirect_uri", redirectURI)
	query.Add("response_mode", "query")
	// TODO: scope from config
	query.Add("scope", "openid profile email")
	// TODO: state and nonce from config/generated
	query.Add("state", "12345")
	query.Add("nonce", "678910")

	req.URL.RawQuery = query.Encode()

	// Redirect the user to the Azure AD login page
	http.Redirect(w, r, req.URL.String(), http.StatusFound)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Authorization code not found", http.StatusBadRequest)
		return
	}
	token, err := exchangeCodeForToken(code)
	if err != nil {
		http.Error(w, "Error exchanging code for token: "+err.Error(), http.StatusUnauthorized)
		return
	}
	// TODO: Encrypt token
	cookie := &http.Cookie{
		Name:     "session",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour), // TODO: Get expire from token
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/api", http.StatusFound)

}

func authProxy(w http.ResponseWriter, r *http.Request) {
	/*
		Will proxy all requests to the backend server, but only if the user is authenticated.
	*/
	cookie, err := r.Cookie("session") // Cookie name from config
	if err != nil {
		// TODO: Redirect to login page
		log.Print("WARNING: Unauthorized request") // TODO: Proper logger
		redirectToAuthenticate(w, r)
		return
	}
	fmt.Fprintf(w, "Cookie found: %s = %s\n", cookie.Name, cookie.Value)
}

func main() {
	http.HandleFunc("/api", authProxy)
	http.HandleFunc("/callback", callbackHandler)
	log.Print("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Print("Error starting server:", err)
	}
}
