package entities

type SlackMetaData struct {
	TeamID         string
	TeamName       string
	WorkspaceToken string
}

// Structs to decode the full Slack server credentials payload missing from slack-go
type SlackCredentials struct {
	ClientID          string `json:"client_id"`
	ClientSecret      string `json:"client_secret"`
	VerificationToken string `json:"verification_token"`
	SigningSecret     string `json:"signing_secret"`
}

type FullManifestResponse struct {
	Ok                bool             `json:"ok"`
	Error             string           `json:"error,omitempty"`
	AppID             string           `json:"app_id,omitempty"`
	Credentials       SlackCredentials `json:"credentials,omitempty"`
	OAuthAuthorizeURL string           `json:"oauth_authorize_url,omitempty"`
}
