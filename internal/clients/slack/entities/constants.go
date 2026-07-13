package entities

import "fmt"

const (
	BOT_SCOPES  string = "users:read"
	USER_SCOPES string = "users.profile:write,users:write,dnd:write"
)

const (
	SLACK_MAX_TEXT_LENGTH int = 100
	SLACK_AUTHORIZER_PORT int = 8080
)

var (
	CALLBACK_URL string = fmt.Sprintf("https://localhost:%d/oauth/callback", SLACK_AUTHORIZER_PORT)
)
