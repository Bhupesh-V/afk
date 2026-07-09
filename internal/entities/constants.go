package entities

// All supported providers
const (
	PROVIDER_SLACK string = "slack"
)

// Names of all AFK secret "keys" being stored in user's system keyring (other than workspace tokens)
const (
	KEYRING_APP_NAME            string = "afk"
	KEYRING_SLACK_CLIENT_ID     string = "afk_slack_app_client_id"
	KEYRING_SLACK_CLIENT_SECRET string = "afk_slack_app_client_secret"
	KEYRING_SLACK_TEAM_IDS      string = "afk_slack_team_ids"
	KEYRING_TEAM_SECRET_PREFIX  string = "afk_slack_team_secret"
	KEYRING_SECRET_SEP          string = "|"
)

// AFK config file constants
const (
	AFK_CONFIG_DIR_NAME string = "afk"
	AFK_CONFIG_FILENAME string = "config.toml"
)

// Slack specific constants
const (
	SLACK_MAX_TEXT_LENGTH int = 100
	SLACK_AUTHORIZER_PORT int = 8080
)
