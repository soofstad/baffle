package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type AccessToken struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
	IdToken      string `json:"id_token"`
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
	var json_token AccessToken
	err = json.Unmarshal([]byte(token), &json_token)
	if err != nil {
		http.Error(w, "Error parsing token: "+err.Error(), http.StatusUnauthorized)
		return
	}
	// TODO: Encrypt token
	cookie := &http.Cookie{
		Name:     "session",
		Value:    token,
		Expires:  time.Now().Add(time.Second * time.Duration(json_token.ExpiresIn)),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	// TODO: Redirect to original request
	http.Redirect(w, r, "/api", http.StatusFound)

}

func authProxy(w http.ResponseWriter, r *http.Request) {
	/*
		Will proxy all requests to the backend server, but only if the user is authenticated.
	*/
	_, err := r.Cookie("session") // Cookie name from config
	if err != nil {
		log.Print("WARNING: Unauthorized request")
		redirectToAuthenticate(w, r)
		return
	}
	// TODO: Validate cookie
	log.Print("INFO: User is authenticated")
	log.Printf("INFO: Proxying request to %s", r.URL.String())
	// TODO: Set auth header and copy all other headers
	proxyRequest, err := http.NewRequest(r.Method, "http://127.0.0.1:80"+r.URL.String(), r.Body)
	resp, err := http.DefaultTransport.RoundTrip(proxyRequest)
	if err != nil {
		log.Print("ERROR: Error creating proxy request: ", err)
		http.Error(w, "Failed to proxy request", 500)
		return
	}

	defer resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {
	http.HandleFunc("/", authProxy)
	http.HandleFunc("/callback", callbackHandler)
	config = loadConfig()
	log.Println("Config loaded: OK")
	log.Print("Starting server on port '8080'")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Print("Error starting server:", err)
	}
}
