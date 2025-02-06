package main

type AccessToken struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
	IdToken      string `json:"id_token"`
}

type Config struct {
	ClientSecret           string
	ClientID               string
	TokenEndpoint          string
	AuthenticationEndpoint string
	RedirectURI            string
	Scope                  string
	CookieName             string
	PathBackendMapping     map[string]string
}
