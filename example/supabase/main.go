package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-zoox/oauth2"
	"github.com/go-zoox/oauth2/supabase"
)

func main() {
	// Environment variables needed:
	// SUPABASE_BASE_URL=https://your-project.supabase.co
	// SUPABASE_CLIENT_ID=your-client-id
	// SUPABASE_CLIENT_SECRET=your-client-secret
	// SUPABASE_REDIRECT_URI=http://localhost:8080/auth/callback

	baseURL := os.Getenv("SUPABASE_BASE_URL")
	clientID := os.Getenv("SUPABASE_CLIENT_ID")
	clientSecret := os.Getenv("SUPABASE_CLIENT_SECRET")
	redirectURI := os.Getenv("SUPABASE_REDIRECT_URI")

	if baseURL == "" || clientID == "" || clientSecret == "" || redirectURI == "" {
		log.Fatal("Missing required environment variables: SUPABASE_BASE_URL, SUPABASE_CLIENT_ID, SUPABASE_CLIENT_SECRET, SUPABASE_REDIRECT_URI")
	}

	// Create Supabase OAuth2 client
	client, err := supabase.New(&supabase.SupabaseConfig{
		BaseURL:      baseURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		Scope:        "openid email profile",
	})
	if err != nil {
		log.Fatal("Failed to create Supabase client:", err)
	}

	// Home page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Supabase OAuth2 Example</title>
			<style>
				body { font-family: Arial, sans-serif; margin: 50px; }
				.button { 
					background-color: #4CAF50; 
					color: white; 
					padding: 14px 20px; 
					text-decoration: none; 
					border: none; 
					border-radius: 4px; 
					cursor: pointer; 
					display: inline-block; 
					margin: 10px 0;
				}
				.button:hover { background-color: #45a049; }
			</style>
		</head>
		<body>
			<h1>Supabase OAuth2 Example</h1>
			<p>This is a simple example of using Supabase OAuth2 authentication.</p>
			<a href="/login" class="button">Login with Supabase</a>
			<a href="/logout" class="button">Logout</a>
		</body>
		</html>
		`
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	})

	// Login endpoint
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		client.Authorize("oauth2-state", func(loginURL string) {
			http.Redirect(w, r, loginURL, http.StatusFound)
		})
	})

	// Logout endpoint
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		client.Logout(func(logoutURL string) {
			if logoutURL == "" {
				// If no logout URL provided, redirect to home
				http.Redirect(w, r, "/", http.StatusFound)
			} else {
				http.Redirect(w, r, logoutURL, http.StatusFound)
			}
		})
	})

	// OAuth callback endpoint
	http.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		client.Callback(code, state, func(user *oauth2.User, token *oauth2.Token, err error) {
			if err != nil {
				log.Printf("OAuth callback error: %v", err)
				http.Error(w, "Authentication failed: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// Successful authentication
			log.Printf("User authenticated: %+v", user)
			log.Printf("Token: %+v", token)

			// Create a simple success page
			html := `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Authentication Success</title>
				<style>
					body { font-family: Arial, sans-serif; margin: 50px; }
					.success { background-color: #d4edda; color: #155724; padding: 20px; border-radius: 4px; margin: 20px 0; }
					.user-info { background-color: #f8f9fa; padding: 20px; border-radius: 4px; margin: 20px 0; }
					.button { 
						background-color: #007bff; 
						color: white; 
						padding: 10px 20px; 
						text-decoration: none; 
						border-radius: 4px; 
						display: inline-block; 
						margin: 10px 0;
					}
				</style>
			</head>
			<body>
				<h1>Authentication Success!</h1>
				<div class="success">
					<h3>Welcome, you have successfully authenticated with Supabase!</h3>
				</div>
				<div class="user-info">
					<h4>User Information:</h4>
					<p><strong>ID:</strong> %s</p>
					<p><strong>Email:</strong> %s</p>
					<p><strong>Username:</strong> %s</p>
					<p><strong>Nickname:</strong> %s</p>
				</div>
				<a href="/" class="button">Go Home</a>
			</body>
			</html>
			`
			response := fmt.Sprintf(html, user.ID, user.Email, user.Username, user.Nickname)
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(response))
		})
	})

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Printf("Visit http://localhost:%s to test the Supabase OAuth2 integration", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}