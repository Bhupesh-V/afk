package slack

import (
	"fmt"
	"log"
	"sync"

	"github.com/slack-go/slack"
)

type Slack interface {
	SetUserCustomStatus(tokens map[string]string, text, emoji string, expiration int64) error
	ClearStatus(tokens map[string]string) error
}

type slackDep struct {
	// Future Deps
}

func New() Slack {
	return &slackDep{}
}

func (s *slackDep) SetUserCustomStatus(tokens map[string]string, text, emoji string, expiration int64) error {
	var wg sync.WaitGroup

	for team, token := range tokens {
		teamName := team
		slackToken := token

		wg.Go(func() {
			api := slack.New(slackToken)

			err := api.SetUserCustomStatus(text, emoji, expiration)
			if err != nil {
				log.Printf("Failed to update status for team %s...: %v", teamName, err)
				return // Exits this specific concurrent worker task cleanly
			}
			fmt.Printf("Successfully updated status using team %s...\n", teamName)
		})
	}

	wg.Wait()

	return nil
}

func (s *slackDep) ClearStatus(tokens map[string]string) error {
	var wg sync.WaitGroup

	for team, token := range tokens {
		teamName := team
		slackToken := token

		wg.Go(func() {
			api := slack.New(slackToken)

			err := api.SetUserCustomStatus("", "", 0)
			if err != nil {
				log.Printf("Failed to update status for team %s...: %v", teamName, err)
				return // Exits this specific concurrent worker task cleanly
			}
			fmt.Printf("Successfully updated status using team %s...\n", teamName)
		})
	}

	wg.Wait()

	return nil
}
