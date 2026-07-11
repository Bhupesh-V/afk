package afk

import (
	"afk/internal/clients/slack"
	"afk/internal/config"
	"afk/internal/entities"
	"afk/pkg"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"
)

type Afk interface {
	// Update user status across configured workspaces
	UpdateStatus(preset, duration string) error
	// Clear user status across configured workspaces
	ClearStatus() error
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

func (a *afk) UpdateStatus(preset, presetDuration string) error {
	var text, emoji, duration string
	if val, ok := a.config.Presets[preset]; ok {
		weekday := time.Now().Weekday()

		// pick emoji and text from 'schedule' if present
		if len(val.Schedule) > 0 {
			for _, item := range val.Schedule {
				if strings.EqualFold(weekday.String(), item.Day) {
					text = item.Text
					if len(item.Emojis) > 0 {
						emoji = random(item.Emojis)
					}
					break
				}
			}
		}

		if presetDuration != "" {
			// priority to CLI value
			duration = presetDuration
		} else {
			if val.Duration != "" {
				duration = val.Duration
			} else {
				fmt.Printf("No duration supplied in config or CLI for preset %s\n", preset)
			}
		}

		// pick emoji from default preset config
		if emoji == "" {
			if len(val.Emojis) > 0 {
				emoji = random(val.Emojis)

			}
		}

		text = val.Text

	} else {
		fmt.Printf("preset '%s' not found, winging it!", preset)
		// TODO create new one by starting a questionaire
	}

	durationUnix, err := pkg.ParseDuration(duration)
	if err != nil {
		return err
	}

	switch a.config.User.Provider {
	case entities.PROVIDER_SLACK:
		err = a.slack.SetUserCustomStatus(
			text,
			fmt.Sprintf(":%s:", emoji),
			time.Now().Add(durationUnix).Unix(),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *afk) ClearStatus() error {
	switch a.config.User.Provider {
	case entities.PROVIDER_SLACK:
		return a.slack.ClearStatus()
	}
	return nil
}

func random(items []string) string {
	if len(items) == 0 {
		return ""
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(items))))
	if err != nil {
		// crypto/rand should rarely fail unless the system entropy pool is broken
		panic(err)
	}

	return items[n.Int64()]
}
