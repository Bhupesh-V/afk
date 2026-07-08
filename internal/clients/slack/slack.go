package slack

import (
	"afk/internal/entities"
	"afk/internal/secrets"
	"afk/pkg/security"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/slack-go/slack"
)

type Slack interface {
	Authorize(clientID, clientSecret string) (string, error)
	SetUserCustomStatus(text, emoji string, expiration int64) error
	ClearStatus() error
	SetPresence(presence string) error
	ToggleNotifications(enabled bool, minutes int) error
}

type slackDep struct {
	uwTokens map[string]string
	secrets  secrets.Secret
}

func New(secrets secrets.Secret) Slack {
	tokens, err := secrets.GetTokens(entities.PROVIDER_SLACK)
	if err != nil {
		return nil
	}

	return &slackDep{
		uwTokens: tokens,
		secrets:  secrets,
	}
}

func (s *slackDep) SetUserCustomStatus(text, emoji string, expiration int64) error {
	var wg sync.WaitGroup

	for team, token := range s.uwTokens {
		teamName := team
		slackToken := token

		wg.Go(func() {
			api := slack.New(slackToken)

			err := api.SetUserCustomStatus(text, emoji, expiration)
			if err != nil {
				log.Printf("Failed to update status for team %s...: %v", teamName, err)
				return // Exits this specific concurrent worker task cleanly
			}
		})
	}

	wg.Wait()

	return nil
}

func (s *slackDep) ClearStatus() error {
	var wg sync.WaitGroup

	for team, token := range s.uwTokens {
		teamName := team
		slackToken := token

		wg.Go(func() {
			api := slack.New(slackToken)

			err := api.SetUserCustomStatus("", "", 0)
			if err != nil {
				log.Printf("Failed to update status for team %s...: %v", teamName, err)
				return // Exits this specific concurrent worker task cleanly
			}
		})
	}

	wg.Wait()

	return nil
}

// SetPresence changes visibility presence. Valid values are "auto" or "away".
func (s *slackDep) SetPresence(presence string) error {
	var wg sync.WaitGroup

	for team, token := range s.uwTokens {
		teamName := team
		slackToken := token

		wg.Go(func() {
			api := slack.New(slackToken)

			err := api.SetUserPresence(presence)
			if err != nil {
				log.Printf("Failed to update presence to %s for team %s: %v", presence, teamName, err)
				return
			}
		})
	}

	wg.Wait()
	return nil
}

// ToggleNotifications handles turning DND off (notifications ON)
func (s *slackDep) ToggleNotifications(enabled bool, minutes int) error {
	var wg sync.WaitGroup

	for team, token := range s.uwTokens {
		teamName := team
		slackToken := token

		wg.Go(func() {
			api := slack.New(slackToken)

			var err error

			if minutes > 0 {
				// snooze until minutes
				_, err = api.SetSnooze(minutes)
			} else {
				_, err = api.EndSnooze()
			}

			if err != nil {
				log.Printf("Failed to toggle notifications for team %s: %v", teamName, err)
				return
			}
		})
	}

	wg.Wait()
	return nil
}

func (s *slackDep) Authorize(clientID, clientSecret string) (string, error) {
	redirectURI := "https://localhost:8080/oauth/callback"
	scopes := "users.profile:write users:write dnd:write"

	base, _ := url.Parse("https://slack.com/oauth/v2/authorize")
	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("user_scope", scopes)
	params.Add("redirect_uri", redirectURI)
	base.RawQuery = params.Encode()

	fmt.Println("--- COPY AND PASTE THIS LINK INTO YOUR BROWSER (where you are logged in from slack) ---")
	fmt.Println(base.String())
	fmt.Println("--------------------------------------------------")

	// Generate the certificate
	tlsCert, err := security.GenerateCertificate()
	if err != nil {
		return "", fmt.Errorf("failed to generate TLS certificates: %w", err)
	}

	mux := http.NewServeMux()

	// Create the server object and supply the TLSConfig holding your in-memory certificate
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
		},
	}

	mux.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing code", http.StatusBadRequest)
			return
		}

		resp, err := slack.GetOAuthV2ResponseContext(context.Background(), &http.Client{}, clientID, clientSecret, code, redirectURI)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("\n--- SUCCESS ---\nWorkspace: %s\nUser Token: %s\n---------------\n", resp.Team.Name, resp.AuthedUser.AccessToken)
		fmt.Fprint(w, "Authorization successful! You can close this tab and check your terminal.")

		log.Println("Shutdown signal received, closing server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Server Shutdown Failed: %+v", err)
		}
	})

	log.Println("Starting server on :8080")

	// Pass empty strings since the certs are already injected via srv.TLSConfig
	if err := srv.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe error: %v", err)
	}

	return "", nil
}
