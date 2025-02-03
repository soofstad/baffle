package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

var config *Config

type Config struct {
	ClientSecret           string
	ClientID               string
	TokenEndpoint          string
	AuthenticationEndpoint string
	RedirectURI            string
	Scope                  string
	CookieName             string
}

func getEnv[T string | int64](key string, defaultValue ...string) T {
	var result T
	var value string
	v, exists := os.LookupEnv(key)
	if !exists {
		if len(defaultValue) > 0 {
			value = defaultValue[0]
		} else {
			log.Fatal("No ", key, " in env")
		}
	} else {
		value = strings.Trim(v, "\"")
	}
	switch any(result).(type) {
	case int64:
		num, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Fatal("Error parsing int: ", err)
		}
		return any(num).(T)
	case string:
		return any(value).(T)
	}
	return any("").(T)
}

func loadConfig() *Config {
	config := &Config{
		ClientSecret:           getEnv[string]("CLIENT_SECRET"),
		ClientID:               getEnv[string]("CLIENT_ID"),
		TokenEndpoint:          getEnv[string]("TOKEN_ENDPOINT"),
		AuthenticationEndpoint: getEnv[string]("AUTHENTICATION_ENDPOINT"),
		RedirectURI:            getEnv[string]("REDIRECT_URI"),
		Scope:                  getEnv[string]("SCOPE"),
		CookieName:             getEnv[string]("COOKIE_NAME", "session"),
	}
	return config
}
