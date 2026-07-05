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
	tokens, err := a.getTokens()
	if err != nil || tokens == nil {
		// everything falls apart if we can't interact with keychain
		return fmt.Errorf("couldn't get secrets from system's keychain")
	}

	var text, emoji, duration string
	if val, ok := a.config.Presets[preset]; ok {
		weekday := time.Now().Weekday()

		if len(val.Schedule) > 0 {
			for _, item := range val.Schedule {
				if strings.EqualFold(weekday.String(), item.Day) {
					text = item.Text
					emoji = random(item.Emojis)
					duration = maxDuration(item.Durations)
					break
				}
			}
		}

		if presetDuration != "" {
			// priority to CLI value
			duration = presetDuration
		} else {
			if duration == "" {
				// find from default preset config
				for _, val := range a.config.Presets {
					if len(val.Durations) > 0 {
						duration = maxDuration(val.Durations)
						break
					} else {
						fmt.Printf("No duration supplied in config or CLI for preset %s\n", preset)
					}
				}
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

	durationUnix, err := pkg.GetTimeDurationFromRelativeDate(duration)
	if err != nil {
		return err
	}

	err = slack.New().SetUserCustomStatus(
		tokens,
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
	tokens, err := a.getTokens()
	if err != nil {
		return err
	}

	return a.slack.ClearStatus(tokens)
}

func (a *afk) getTokens() (map[string]string, error) {
	// TODO: figure out tokens
	return nil, nil
}

func random(items []string) string {
	return items[rand.IntN(len(items))]
}

func maxDuration(items []string) string {
	max, err := pkg.GetMaxDurationFromRelativeDates(items)
	if err != nil {
		return ""
	}

	return max.String()
}
