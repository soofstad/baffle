# GoathBFF
A simple, fast, and configurable OAuth2 Backend-for-frontend

## DISCLAIMER
- This is a work-in-progress and should __not be used in production__
- Only meant to support single, stateless, API-calls. Not meant for serving HTML pages or static assets
- Probably not very secure or reliable at the moment... but we're working on it!
- Also not the greatest of all time (yet)

## Example

This example uses docker-compose to run a simple setup with a backend and the GoathBFF.  
The GoathBFF will proxy requests to the backend and handle OAuth2 authentication.  
The backend will receive a regular "Authorization" header with the access/bearer token.

Your frontend should be using GoathBFF for all API requests.  
Other clients (confidential clients) should obtain an access token in other ways, and send requests directly to the backend. 

```bash

```yaml
services:
  goathbff:
    image: soofstad/goathbff:latest
    environment:
      # OAuth2 configuration
      CLIENT_SECRET: ${CLIENT_SECRET}
      CLIENT_ID: 5cb6c4de-28d0-4b62-a547-262dc2377baf
      REDIRECT_URI: http://localhost:8080/callback
      TOKEN_ENDPOINT: "https://login.microsoftonline.com/3aa4a235-b6e2-48d5-9195-7fcf05b459b0/oauth2/v2.0/token"
      AUTHENTICATION_ENDPOINT: "https://login.microsoftonline.com/3aa4a235-b6e2-48d5-9195-7fcf05b459b0/oauth2/v2.0/authorize"
      SCOPE: "openid profile email"
      
      # Paths that should be proxied to the backend(s)
      # Note that '/callback' and '/session' are reserved by baffle
      # Format: "path,backend-url;path,backend-url" (empty path is allowed)
      PATH_BACKEND_MAPPING: "api,http://backend;,http://backend"
      
      # Optional - will default to "session"
      COOKIE_NAME: "session"
    ports:
      - "80:8080"
  
  backend:
    image: nginx
    ports:
      - "8080:80"
```

## How-to
TODO: Something about correct IdentityProvider setup. Do not allow SPA or implicit flow for example...