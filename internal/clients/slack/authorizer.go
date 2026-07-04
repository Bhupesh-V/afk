package slack

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/slack-go/slack"
)

func Authorize(clientID, clientSecret string) (string, error) {
	redirectURI := "https://localhost:8080/oauth/callback"
	scopes := "users.profile:write"

	// Dynamically build a flawless, URL-encoded string
	base, _ := url.Parse("https://slack.com/oauth/v2/authorize")
	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("user_scope", scopes)
	params.Add("redirect_uri", redirectURI)
	base.RawQuery = params.Encode()

	fmt.Println("--- COPY AND PASTE THIS LINK INTO YOUR BROWSER ---")
	fmt.Println(base.String())
	fmt.Println("--------------------------------------------------")

	mux := http.NewServeMux()

	// Create the server object explicitly
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing code", http.StatusBadRequest)
			return
		}

		// Exchange the temporary code for a permanent User Token
		resp, err := slack.GetOAuthV2ResponseContext(context.Background(), &http.Client{}, clientID, clientSecret, code, redirectURI)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Print your permanent user token to the console
		fmt.Printf("\n--- SUCCESS ---\nWorkspace: %s\nUser Token: %s\n---------------\n", resp.Team.Name, resp.AuthedUser.AccessToken)
		fmt.Fprint(w, "Authorization successful! You can close this tab and check your terminal.")

		log.Println("Shutdown signal received, closing server...")

		// Establish a timeout window for other active requests to finish
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Server Shutdown Failed: %+v", err)
		}
	})

	log.Println("Starting server on :8080")

	// ListenAndServe blocks until the server is stopped
	if err := srv.ListenAndServeTLS("server.pem", "server.key"); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe error: %v", err)
	}

	return "", fmt.Errorf("failed to start OAuth server")
}
