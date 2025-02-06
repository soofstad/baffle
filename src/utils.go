package main

import (
	"errors"
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
	HostWhiteListPaths     map[string]string
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
	rawHostWhiteListPaths := getEnv[string]("HOST_WHITELIST_PATHS", "")
	if rawHostWhiteListPaths == "" {
		log.Fatal("No whitelisted backend hosts found. The BFF can't possibly work.")
	}
	// TODO: Do some sanity checks on the strings
	pairs := strings.Split(rawHostWhiteListPaths, ";")
	hostWhiteListPaths := make(map[string]string)
	for _, pair := range pairs {
		parts := strings.Split(pair, ",")
		if len(parts) != 2 {
			log.Fatalf("Invalid HOST_WHITELIST_PATHS %s. Should be on format 'aliasA,backendA;aliasB,backendB", pair)
		}
		hostWhiteListPaths[parts[0]] = parts[1]
	}
	config := &Config{
		ClientSecret:           getEnv[string]("CLIENT_SECRET"),
		ClientID:               getEnv[string]("CLIENT_ID"),
		TokenEndpoint:          getEnv[string]("TOKEN_ENDPOINT"),
		AuthenticationEndpoint: getEnv[string]("AUTHENTICATION_ENDPOINT"),
		RedirectURI:            getEnv[string]("REDIRECT_URI"),
		Scope:                  getEnv[string]("SCOPE"),
		CookieName:             getEnv[string]("COOKIE_NAME", "session"),
		HostWhiteListPaths:     hostWhiteListPaths,
	}
	return config
}

func getProxyTargetFromPath(path string) (string, error) {
	trimmedPath := strings.Trim(path, "/")
	pathParts := strings.SplitN(trimmedPath, "/", 2)
	alias := pathParts[0]
	backend, exists := config.HostWhiteListPaths[alias]
	if !exists {
		log.Print("ERROR: No configured backend found for alias" + alias)
		return "", errors.New("No configured backend found for alias " + alias)
	}
	targetPath := strings.Join(pathParts[1:], "/")
	backend = strings.Trim(backend, "/")
	return backend + "/" + targetPath, nil
}
