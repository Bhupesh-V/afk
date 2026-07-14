package slack

import (
	slackEntities "afk/internal/clients/slack/entities"
	"afk/internal/clients/slack/errors"
	"afk/internal/entities"
	"afk/internal/secrets"
	"afk/pkg/security"
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"

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
	ClearStatus() error
	DispatchStatus(opts ...Option) error
}

type slackDep struct {
	uwTokens map[string]string
	secrets  secrets.Secret
}

// statusConfig holds the parameters gathered from the options
type statusConfig struct {
	text          string
	emoji         string
	expiration    int64
	presence      string
	notifications bool
}

type Option func(*statusConfig)

// WithStatus sets the custom status text, emoji, and expiration
func WithStatus(text, emoji string, durationUnix time.Duration) Option {
	return func(c *statusConfig) {
		c.text = text
		c.emoji = fmt.Sprintf(":%s:", emoji)
		c.expiration = time.Now().Add(durationUnix).Unix()
	}
}

// WithPresence sets the user presence
func WithPresence(presence string) Option {
	return func(c *statusConfig) {
		c.presence = presence
	}
}

// WithNotifications sets DND
func WithNotifications(dnd bool) Option {
	return func(c *statusConfig) {
		c.notifications = dnd
	}
}

func New(secrets secrets.Secret, isSetup bool) (Slack, error) {
	deps := &slackDep{
		secrets: secrets,
	}

	if isSetup {
		// TODO Ask user to build a config

		clientId, clientSecret, err := deps.createSlackApp()
		if err != nil {
			return nil, fmt.Errorf("failed to create slack app, please try again")
		}

		if clientId != "" && clientSecret != "" {
			err := deps.updateAppCreds(clientId, clientSecret)
			if err != nil {
				return nil, err
			}

			fmt.Println("\nLet's install afk to your preferred workspace(s), please follow on-screen instructions")

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
		} else {
			return nil, fmt.Errorf("malformed slack app creds, please try again")
		}
	}

	tokens, err := deps.getTokens()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch slack tokens: %w", err)
	}

	deps.uwTokens = tokens

	return deps, nil
}

func (s *slackDep) Authorize(clientID, clientSecret string) (*slackEntities.SlackMetaData, error) {
	metadata := &slackEntities.SlackMetaData{}

	base, _ := url.Parse("https://slack.com/oauth/v2/authorize")
	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("scope", slackEntities.BOT_SCOPES)
	params.Add("user_scope", slackEntities.USER_SCOPES)

	base.RawQuery = params.Encode()

	fmt.Println()
	fmt.Println("--- OPEN THIS LINK IN YOUR BROWSER (where you are logged in to slack) ---")
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
		Addr:    fmt.Sprintf(":%d", slackEntities.SLACK_AUTHORIZER_PORT),
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

		resp, err := slack.GetOAuthV2ResponseContext(context.Background(), &http.Client{}, clientID, clientSecret, code, slackEntities.CALLBACK_URL)
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
		return nil, fmt.Errorf("error serving oauth callback server %s", err)
	}

	return metadata, nil
}

func (s *slackDep) DispatchStatus(opts ...Option) error {
	var wg sync.WaitGroup
	var cfg statusConfig

	for _, opt := range opts {
		opt(&cfg)
	}

	var mu sync.Mutex
	collectedErrors := make(errors.TeamErrors)

	// Helper function to safely log errors per team across goroutines
	addError := func(team string, err error) {
		mu.Lock()
		defer mu.Unlock()
		collectedErrors[team] = append(collectedErrors[team], err)
	}

	for team, token := range s.uwTokens {
		teamName := team
		slackToken := token

		wg.Go(func() {
			api := slack.New(slackToken)

			var innerWg sync.WaitGroup

			if cfg.emoji != "" || cfg.text != "" {
				innerWg.Go(func() {
					err := api.SetUserCustomStatus(cfg.text, cfg.emoji, cfg.expiration)
					if err != nil {
						addError(teamName, fmt.Errorf("status update failed: %w", err))
						return
					}
				})
			}

			if cfg.presence != "" {
				innerWg.Go(func() {
					err := api.SetUserPresence(cfg.presence)
					if err != nil {
						addError(teamName, fmt.Errorf("presence update failed: %w", err))
						return
					}
				})
			}

			if !cfg.notifications {
				innerWg.Go(func() {
					// send in minutes
					_, err := api.SetSnooze(int(cfg.expiration))
					if err != nil {
						addError(teamName, fmt.Errorf("snooze update failed: %w", err))
						return
					}
				})
			}
			innerWg.Wait()

			mu.Lock()
			_, hasErrors := collectedErrors[teamName]
			mu.Unlock()

			if !hasErrors {
				fmt.Printf("Status updated on %s\n", teamName)
			}
		})
	}
	wg.Wait()

	if len(collectedErrors) > 0 {
		return collectedErrors
	}

	return nil
}

func (s *slackDep) ClearStatus() error {
	var wg sync.WaitGroup

	for team, token := range s.uwTokens {
		teamName := team
		slackToken := token

		wg.Go(func() {
			api := slack.New(slackToken)

			var innerWg sync.WaitGroup

			innerWg.Go(func() {
				err := api.SetUserCustomStatus("", "", 0)
				if err != nil {
					fmt.Printf("Failed to update status for team %s...: %v", teamName, err)
					return
				}
			})

			innerWg.Go(func() {
				err := api.SetUserPresence("auto")
				if err != nil {
					log.Printf("Failed to update presence for team %s: %v", teamName, err)
					return
				}
			})

			innerWg.Go(func() {
				_, err := api.EndSnooze()
				if err != nil {
					return
				}
			})
			innerWg.Wait()

			fmt.Printf("Status updated on %s\n", teamName)
		})
	}
	wg.Wait()

	return nil
}

func (s *slackDep) createSlackApp() (clientId, clientSecret string, err error) {
	clientId, clientSecret, err = s.getAppCreds()
	if err != nil {
		return "", "", err
	}
	if clientId != "" && clientSecret != "" {
		return clientId, clientSecret, nil
	}

	var appConfigurationToken string

	fmt.Println("Looks like this is your first AFK run, running setup wizard")
	fmt.Println()
	fmt.Print("Enter Slack App Configuration Access Token: ")
	fmt.Scanln(&appConfigurationToken)

	manifestStruct := &slack.Manifest{
		Display: slack.Display{
			Name:        "AFK",
			Description: "App used to mark a user availability",
		},
		Settings: slack.Settings{},
		Features: slack.Features{
			BotUser: slack.BotUser{
				DisplayName:  "AFK Bot",
				AlwaysOnline: false,
			},
		},
		OAuthConfig: slack.OAuthConfig{
			RedirectUrls: []string{
				slackEntities.CALLBACK_URL,
			},
			Scopes: slack.OAuthScopes{
				Bot:  strings.Split(slackEntities.BOT_SCOPES, ","),
				User: strings.Split(slackEntities.USER_SCOPES, ","),
			},
		},
	}
	apiURL := "https://slack.com/api/apps.manifest.create"

	manifestBytes, err := json.Marshal(manifestStruct)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal manifest layout: %v", err)
	}
	manifestJSON := string(manifestBytes)

	// Encode parameters as application/x-www-form-urlencoded
	data := url.Values{}
	data.Set("token", appConfigurationToken)
	data.Set("manifest", manifestJSON)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", "", fmt.Errorf("failed to build http request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("http communication error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read stream payload: %v", err)
	}

	var fullResp slackEntities.FullManifestResponse
	if err := json.Unmarshal(body, &fullResp); err != nil {
		return "", "", fmt.Errorf("failed to parse credential fields: %v. Raw body: %s", err, string(body))
	}

	if !fullResp.Ok {
		return "", "", fmt.Errorf("slack error: %s (Raw body: %s)", fullResp.Error, string(body))
	}

	// fmt.Println("Multi-workspace setup:")
	// fmt.Printf("Open the app: https://api.slack.com/apps/%s/general\n", fullResp.AppID)
	// fmt.Println("Head over to 'Manage Distribution' and enable 'Public Distribution'")

	return fullResp.Credentials.ClientID, fullResp.Credentials.ClientSecret, nil
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

	for id := range strings.SplitSeq(rawTeamIds, ",") {
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

	return tokens, nil
}

func (s *slackDep) updateTokens(meta *slackEntities.SlackMetaData) error {
	// TODO fix for new setup
	rawTeamIds, err := s.secrets.GetSecret(entities.KEYRING_SLACK_TEAM_IDS)
	if err != nil {
		return err
	}

	teamIds := strings.Split(rawTeamIds, ",")

	if rawTeamIds == "" {
		teamIds = []string{}
	}

	var found bool
	for _, id := range teamIds {
		if id == meta.TeamID {
			break
		}
	}

	if !found {
		// save the team id in the team list secret
		teamIds = append(teamIds, meta.TeamID)
		err = s.secrets.SetSecret(entities.KEYRING_SLACK_TEAM_IDS, strings.Join(teamIds, ","))
		if err != nil {
			return err
		}
	}

	tokenKey := fmt.Sprintf("%s_%s", entities.KEYRING_TEAM_SECRET_PREFIX, meta.TeamID)
	// overwrite/create the secret
	err = s.secrets.SetSecret(tokenKey, fmt.Sprintf("%s|%s", meta.TeamName, meta.WorkspaceToken))
	if err != nil {
		return err
	}

	return nil
}
