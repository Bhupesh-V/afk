package slack

import (
	"log"
	"sync"

	"github.com/slack-go/slack"
)

type Slack interface {
	SetUserCustomStatus(text, emoji string, expiration int64) error
	ClearStatus() error
	SetPresence(presence string) error
	ToggleNotifications(enabled bool, minutes int) error
}

type slackDep struct {
	uwTokens map[string]string
}

func New(tokens map[string]string) Slack {
	return &slackDep{
		uwTokens: tokens,
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
