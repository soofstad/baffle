package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", proxyHandler)
	http.HandleFunc("/callback", callbackHandler)
	// TODO: This
	//http.HandleFunc("/session", sessionInfo)
	config = loadConfig()
	log.Println("Config loaded: OK")
	log.Print("Backend mapping: ", config.PathBackendMapping)
	log.Print("Starting server on port '8080'")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Print("Error starting server:", err)
	}
}
