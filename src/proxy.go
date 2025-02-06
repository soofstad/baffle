package main

import (
	"io"
	"log"
	"net/http"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	/*
		Will proxy all requests to the backend server matching the first path element.
		Requires that the request has a cookie named the same as configured in 'config.CookieName'
		The cookie must be decryptable given the configured 'config.GoathBFFSecret'.
		If the access_token has expired we will try to get a new one using the refresh_token.
		TODO: ?If the refresh_token has expired, the request will fail with a 401.?
		There is no checks or validation on the access token itself.
	*/
	sessionCookie, err := r.Cookie(config.CookieName)
	// TODO: Decrypt cooke to get token
	if err != nil {
		log.Printf("IFNO: No cookie named '%s' found. Redirecting to authenticate", config.CookieName)
		redirectToAuthenticate(w, r, r.URL.String())
		return
	}
	log.Print("INFO: User is authenticated")
	backend, err := getProxyTargetFromPath(r.URL.String())
	if err != nil {
		http.Error(w, "Error getting backend: "+err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("INFO: Proxying request to %s", backend)
	proxyRequest, _ := http.NewRequest(r.Method, backend, r.Body)
	// TODO: Do we need to copy headers?
	// TODO: Handle expired access/refresh tokens
	proxyRequest.Header.Set("Authorization", "Bearer "+sessionCookie.Value)
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
