package afk

import (
	"afk/internal/clients/slack"
	"afk/internal/config"
)

type Afk interface {
	UpdateStatus(preset string) error
}

type afk struct {
	config config.Config
	slack  slack.Slack
}

func New(config config.Config, slack slack.Slack) Afk {
	return &afk{
		config: config,
		slack:  slack,
	}
}

func (a *afk) UpdateStatus(preset string) error {
	// TODO get token from system
	// TODO get preset text and emoji from config
	return nil
}
