package afk

import (
	"afk/internal/clients/slack"
	"afk/internal/config"
	"afk/internal/entities"
	"afk/pkg"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/bhupesh-v/promptui"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	emojilib "github.com/mrosales/emoji-go"
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
	var text, emoji, duration, presence string
	var dnd bool = true

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

		// take from default
		text = val.Text

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

		if val.Presence != "" {
			presence = val.Presence
		}

		if val.Notifications == "on" {
			dnd = false
		}

	} else {
		fmt.Printf("Preset '%s' not found, winging it!\n\n", preset)

		// TODO move to interface/cmd

		// Gather Status Text & Duration (standard huh form)
		var firstFields []huh.Field
		firstFields = append(firstFields,
			huh.NewInput().
				Title("Status Text").
				Prompt("» ").
				Placeholder("Going on a walk...").
				Value(&text).
				Validate(func(str string) error {
					if len(str) > 100 {
						return fmt.Errorf("text exceeds max length of %d (got %d chars)", 100, len(str))
					}
					return nil
				}),
		)

		if presetDuration != "" {
			duration = presetDuration
		} else {
			firstFields = append(firstFields,
				huh.NewInput().
					Title("Duration").
					Prompt("» ").
					Placeholder("e.g., 30m, 1h, 2h30m").
					Value(&duration).
					Validate(func(str string) error {
						_, err := pkg.ParseDuration(duration)
						if err != nil {
							return err
						}
						return nil
					}),
			)
		}

		if err := huh.NewForm(huh.NewGroup(firstFields...).WithTheme(entities.MonokaiRetroTheme())).Run(); err != nil {
			return fmt.Errorf("failed to collect status text or duration: %w", err)
		}

		allEmojis := emojilib.All

		var emojiDisplayList []string
		for _, e := range allEmojis {
			emojiDisplayList = append(emojiDisplayList, fmt.Sprintf(" %s : %s", e.Character, e.Name))
		}

		theme := entities.MonokaiRetroTheme()
		titleStyle := lipgloss.NewStyle().Foreground(theme.Focused.Title.GetForeground()).Bold(true)
		arrowStyle := lipgloss.NewStyle().Foreground(theme.Focused.TextInput.Prompt.GetForeground())

		// Style promptui to mimic the Charm Huh design language
		templates := &promptui.SelectTemplates{
			Label:       "{{ . }}",
			Active:      fmt.Sprintf("%s {{ . }}", arrowStyle.Render("»")),
			SearchLabel: titleStyle.Render("Choose Emoji 🔍︎ "),
		}

		prompt := promptui.Select{
			Label:             " ",
			Items:             emojiDisplayList,
			StartInSearchMode: true,
			Size:              10, // Height
			Templates:         templates,
			HideSelected:      true,
			Stdout:            &pkg.NoBellWriter{Writer: os.Stdout},

			Searcher: func(input string, index int) bool {
				item := strings.ToLower(emojiDisplayList[index])
				input = strings.ToLower(input)

				inputIdx := 0
				for i := 0; i < len(item); i++ {
					if inputIdx < len(input) && item[i] == input[inputIdx] {
						inputIdx++
					}
				}
				return inputIdx == len(input)
			},
		}

		// Launch promptui selection screen
		idx, _, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("emoji selection canceled or failed: %w", err)
		}

		emoji = allEmojis[idx].Name
	}

	durationUnix, err := pkg.ParseDuration(duration)
	if err != nil {
		return err
	}

	switch a.config.User.Provider {
	case entities.PROVIDER_SLACK:
		err = a.slack.DispatchStatus(
			slack.WithStatus(text, emoji, durationUnix),
			slack.WithPresence(presence),
			slack.WithNotifications(dnd),
		)
		if err != nil {
			return err
		}
	}

	// save the preset if one doesn't exist
	if _, ok := a.config.Presets[preset]; !ok {
		a.config.Presets[preset] = config.PresetItem{
			Text:     text,
			Emojis:   []string{emoji},
			Duration: duration,
			Presence: presence,
		}
		err = a.config.Write()
		if err != nil {
			return fmt.Errorf("failed to update afk config with preset [%s]: %v", preset, err)
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
