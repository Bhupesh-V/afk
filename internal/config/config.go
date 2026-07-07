package config

import (
	"afk/internal/entities"
	"afk/pkg"
	"errors"
	"fmt"
	"os"

	"github.com/mrosales/emoji-go"
	"github.com/pelletier/go-toml/v2"
)

// TODO: use "comment" struct tag on each field
type Config struct {
	emojiIndex *emoji.SearchIndex
	User       UserSection           `toml:"user"`
	Presets    map[string]PresetItem `toml:"presets"`
}

type UserSection struct {
	DurationStrategy string `toml:"durationStrategy"`
	Presence         bool   `toml:"presence"`
}

type ScheduleItem struct {
	Day    string   `toml:"day"`
	Text   string   `toml:"text,omitempty"`
	Emojis []string `toml:"emojis,omitempty"`
}

type PresetItem struct {
	Text     string         `toml:"text"`
	Emojis   []string       `toml:"emojis,omitempty"`
	Duration string         `toml:"duration,omitempty"`
	Schedule []ScheduleItem `toml:"schedule,omitempty"`
}

// Validate ensures that no text field in defaults or schedules exceeds 100 characters.
func (c *Config) Validate() error {
	var errs []error

	for name, item := range c.Presets {
		textErrors := c.validateText(item, name)
		errs = append(errs, textErrors...)

		emojiErrors := c.validateEmojis(item, name)
		errs = append(errs, emojiErrors...)

		durationErrors := c.validateDurations(item, name)
		errs = append(errs, durationErrors...)
	}

	return errors.Join(errs...)
}

func New() (Config, error) {
	cfgFile, err := pkg.GetConfigPath(entities.AFK_CONFIG_FILENAME)
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("afk config read failed:\n%w", err)
	}

	cfg.emojiIndex = emoji.NewSearchIndex()

	if err := cfg.Validate(); err != nil {
		return Config{}, fmt.Errorf("afk config validation failed:\n%w", err)
	}

	return cfg, nil
}

func (c *Config) validateText(preset PresetItem, name string) []error {
	const maxLen = entities.SLACK_MAX_TEXT_LENGTH
	var errs []error

	if len(preset.Text) > maxLen {
		errs = append(errs, fmt.Errorf("preset [%s] text exceeds max length of %d (got %d chars)", name, maxLen, len(preset.Text)))
	}

	// Check each day-specific schedule text length as well
	for _, schedule := range preset.Schedule {
		if len(schedule.Text) > maxLen {
			errs = append(errs, fmt.Errorf("preset [%s] schedule for %s text exceeds max length of %d (got %d chars)", name, schedule.Day, maxLen, len(schedule.Text)))
		}
	}

	return errs
}

func (c *Config) validateEmojis(preset PresetItem, name string) []error {
	var errs []error

	if len(preset.Emojis) == 0 {
		errs = append(errs, fmt.Errorf("missing default emoji on preset [%s]", name))
	}

	// collect all emojis
	var emojis []string
	emojis = append(emojis, preset.Emojis...)

	for _, sch := range preset.Schedule {
		emojis = append(emojis, sch.Emojis...)
	}

	// search them in our dataset
	for _, e := range emojis {
		info := c.emojiIndex.Search(e, emoji.WithMaxDistance(0), emoji.WithLimit(1))

		exactMatchFound := false
		if len(info) > 0 {
			for _, name := range info[0].AlternateNames {
				if name == e {
					exactMatchFound = true
					break
				}
			}
		}

		if !exactMatchFound {
			errs = append(errs, fmt.Errorf("invalid emoji [%s] on preset [%s]", e, name))
		}
	}

	return errs
}

func (c *Config) validateDurations(preset PresetItem, name string) []error {
	var errs []error

	if preset.Duration != "" {
		if !pkg.IsValidRelativeDate(preset.Duration) {
			errs = append(errs, fmt.Errorf("found invalid duration [%s] on preset [%s]", preset.Duration, name))
		}
	} else {
		errs = append(errs, fmt.Errorf("missing default duration on preset [%s]", name))
	}

	for _, sch := range preset.Schedule {
		if !pkg.IsValidDay(sch.Day) {
			errs = append(errs, fmt.Errorf("found invalid day [%s] on preset [%s]", sch.Day, name))
		}
	}

	return errs
}
