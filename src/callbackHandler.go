package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	/*
		Should receive an authentication CODE in the 'code' parameter.
		1. Will exchange that code for an access_token and refresh_token.
		2. Encrypt the token response and set it in the cookie
		3. Redirect back to the original URL that was requested before the client had a sessionCookie
	*/
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Authorization code not found", http.StatusBadRequest)
		return
	}
	state := r.URL.Query().Get("state")
	if state == "" {
		log.Print("ERROR: State not found. Can't redirect to original request")
		http.Error(w, "An unhandled server error occurred", http.StatusInternalServerError)
		return
	}
	token, err := exchangeCodeForToken(code)
	if err != nil {
		log.Print("ERROR: Error exchanging code for token: ", err)
		http.Error(w, "An unhandled server error occurred", http.StatusInternalServerError)
		return
	}
	// TODO: Handle non-jwt tokens
	var json_token AccessToken
	err = json.Unmarshal([]byte(token), &json_token)
	if err != nil {
		log.Print("ERROR: Error parsing token: ", err)
		http.Error(w, "An unhandled server error occurred", http.StatusInternalServerError)
		return
	}
	// TODO: Encrypt token
	cookie := &http.Cookie{
		Name:     config.CookieName,
		Value:    token,
		Expires:  time.Now().Add(time.Second * time.Duration(json_token.ExpiresIn)),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	// TODO: Redirect to original request
	http.Redirect(w, r, state, http.StatusFound)

}
