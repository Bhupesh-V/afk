package afk

import (
	"afk/internal/clients/slack"
	"afk/internal/config"
	"afk/pkg"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"
)

type Afk interface {
	UpdateStatus(preset, duration string) error
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

		if emoji == "" {
			// take one from default preset config, if present
			for _, val := range a.config.Presets {
				if len(val.Emojis) > 0 {
					emoji = random(val.Emojis)
					break
				}
			}
		}

	} else {
		fmt.Println("preset not found")
		// TODO create new one by starting a questionaire
	}

	durationUnix, err := pkg.ParseDuration(duration)
	if err != nil {
		return err
	}

	err = a.slack.SetUserCustomStatus(
		text,
		fmt.Sprintf(":%s:", emoji),
		time.Now().Add(durationUnix).Unix(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (a *afk) ClearStatus() error {
	return a.slack.ClearStatus()
}

func random(items []string) string {
	return items[rand.IntN(len(items))]
}
