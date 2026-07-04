package entities

// Names of all AFK secret "keys" being stored in user's system keyring (other than workspace tokens)
const (
	KEYRING_APP_NAME                string = "afk"
	KEYRING_AFK_SLACK_CLIENT_ID     string = "afk_slack_app_client_id"
	KEYRING_AFK_SLACK_CLIENT_SECRET string = "afk_slack_app_client_secret"
)

// AFK config file constants
const (
	AFK_CONFIG_DIR_NAME string = "afk"
	AFK_CONFIG_FILENAME string = ".afk.config.toml"
)
