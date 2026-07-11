package slack

import (
	slackEntities "afk/internal/clients/slack/entities"
	"afk/internal/entities"
	"afk/internal/secrets"
	"afk/pkg/security"
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/slack-go/slack"
)

type Slack interface {
	Authorize(clientID, clientSecret string) (*slackEntities.SlackMetaData, error)
	SetUserCustomStatus(text, emoji string, expiration int64) error
	ClearStatus() error
	SetPresence(presence string) error
	ToggleNotifications(enabled bool, minutes int) error
}

type slackDep struct {
	uwTokens map[string]string
	secrets  secrets.Secret
}

func New(secrets secrets.Secret, isSetup bool) (Slack, error) {
	deps := &slackDep{
		secrets: secrets,
	}

	if isSetup {
		// TODO Ask user to build a config

		var clientId, clientSecret string

		fmt.Println()
		fmt.Print("Enter Slack Client Id: ")
		fmt.Scanln(&clientId)

		fmt.Print("Enter Slack Client Secret: ")
		fmt.Scanln(&clientSecret)

		if clientId != "" && clientSecret != "" {
			err := deps.updateAppCreds(clientId, clientSecret)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("malformed slack app creds, please try again")
		}

		fmt.Println("\nLet's install afk to your preferred workspace, please follow on-screen instructions")

		var addMore bool = true
		reader := bufio.NewReader(os.Stdin)

		for addMore {
			metadata, err := deps.Authorize(clientId, clientSecret)
			if err != nil {
				return nil, err
			}
			err = deps.updateTokens(metadata)
			if err != nil {
				return nil, err
			}

			fmt.Println()
			fmt.Print("Add more workspaces? [y/N]: ")

			// Read just the first byte
			char, err := reader.ReadByte()
			if err != nil {
				return nil, fmt.Errorf("failed to read input: %w", err)
			}

			// If they typed extra characters, flush the rest of the line
			// so it doesn't bleed into the next loop iteration
			if char != '\n' {
				_, _ = reader.ReadString('\n')
			}

			// Convert byte to lowercase string for easy comparison
			choice := strings.ToLower(string(char))

			// Default to exit. Only continue if they explicitly typed 'y'
			if choice != "y" {
				addMore = false
			}
		}

	}

	tokens, err := deps.getTokens()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch slack tokens: %w", err)
	}

	deps.uwTokens = tokens

	return deps, nil
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

func (s *slackDep) Authorize(clientID, clientSecret string) (*slackEntities.SlackMetaData, error) {
	metadata := &slackEntities.SlackMetaData{}
	redirectURI := fmt.Sprintf("https://localhost:%d/oauth/callback", entities.SLACK_AUTHORIZER_PORT)
	scopes := "users.profile:write users:write dnd:write"

	base, _ := url.Parse("https://slack.com/oauth/v2/authorize")
	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("user_scope", scopes)
	params.Add("redirect_uri", redirectURI)
	base.RawQuery = params.Encode()

	fmt.Println()
	fmt.Println("--- COPY AND PASTE THIS LINK INTO YOUR BROWSER (where you are logged in from slack) ---")
	fmt.Println(base.String())
	fmt.Println("--------------------------------------------------")

	// Generate the certificate
	tlsCert, err := security.GenerateCertificate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate TLS certificates: %w", err)
	}

	mux := http.NewServeMux()

	// Create the server object and supply the TLSConfig holding your in-memory certificate
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", entities.SLACK_AUTHORIZER_PORT),
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

		fmt.Printf("\n--- SUCCESS ---\nWorkspace: %s\n", resp.Team.Name)
		fmt.Fprint(w, "Authorization successful! You can close this tab and check your terminal.")

		// Force push data out of the buffer and into the network (pretty interesting)
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}

		metadata.TeamID = resp.Team.ID
		metadata.TeamName = resp.Team.Name
		metadata.WorkspaceToken = resp.AuthedUser.AccessToken

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				fmt.Printf("Server Shutdown Failed: %+v\n", err)
			}
		}()
	})

	// Pass empty strings since the certs are already injected via srv.TLSConfig
	if err := srv.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe error: %v", err)
	}

	return metadata, nil
}

func (s *slackDep) getAppCreds() (string, string, error) {
	var clientId, clientSecret string
	var err error

	clientId, err = s.secrets.GetSecret(entities.KEYRING_SLACK_CLIENT_ID)
	if err != nil {
		return clientId, clientSecret, err
	}

	clientSecret, err = s.secrets.GetSecret(entities.KEYRING_SLACK_CLIENT_SECRET)
	if err != nil {
		return clientId, clientSecret, err
	}

	return clientId, clientSecret, nil
}

func (s *slackDep) updateAppCreds(id, secret string) error {
	var err error

	err = s.secrets.SetSecret(entities.KEYRING_SLACK_CLIENT_ID, id)
	if err != nil {
		return err
	}

	err = s.secrets.SetSecret(entities.KEYRING_SLACK_CLIENT_SECRET, secret)
	if err != nil {
		return err
	}

	return nil
}

func (s *slackDep) getTokens() (map[string]string, error) {
	tokens := make(map[string]string)

	rawTeamIds, err := s.secrets.GetSecret(entities.KEYRING_SLACK_TEAM_IDS)
	if err != nil {
		return nil, err
	}

	for _, id := range strings.Split(rawTeamIds, ",") {
		tokenKey := fmt.Sprintf("%s_%s", entities.KEYRING_TEAM_SECRET_PREFIX, id)
		token, err := s.secrets.GetSecret(tokenKey)
		if err != nil {
			return nil, err
		}
		// rawTokens[0] is team namere
		// rawTokens[1] is token
		rawTokens := strings.Split(token, entities.KEYRING_SECRET_SEP)
		if len(rawTokens) != 2 {
			return nil, fmt.Errorf("inconsisten token format on system keyring")
		}
		tokens[rawTokens[0]] = rawTokens[1]
	}

	return nil, nil
}

func (s *slackDep) updateTokens(meta *slackEntities.SlackMetaData) error {
	// TODO fix for new setup
	rawTeamIds, err := s.secrets.GetSecret(entities.KEYRING_SLACK_TEAM_IDS)
	if err != nil {
		return err
	}
	teamIds := strings.Split(rawTeamIds, ",")

	var found bool
	for _, id := range teamIds {
		if id == meta.TeamID {
			break
		}
	}
	tokenKey := fmt.Sprintf("%s_%s", entities.KEYRING_TEAM_SECRET_PREFIX, meta.TeamID)

	if !found {
		// save the team id in the team list secret
		teamIds = append(teamIds, meta.TeamID)
		err = s.secrets.SetSecret(entities.KEYRING_SLACK_TEAM_IDS, strings.Join(teamIds, ","))
		if err != nil {
			return err
		}
	}

	// overwrite/create the secret
	err = s.secrets.SetSecret(tokenKey, fmt.Sprintf("%s|%s", meta.TeamName, meta.WorkspaceToken))
	if err != nil {
		return err
	}

	return nil
}
