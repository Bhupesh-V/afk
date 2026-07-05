package config

import (
	"afk/internal/entities"
	"afk/pkg"
	"errors"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// TODO: use "comment" struct tag on each field
type Config struct {
	User    UserSection           `toml:"user"`
	Presets map[string]PresetItem `toml:"presets"`
}

type UserSection struct {
	DurationStrategy string `toml:"durationStrategy"`
	Presence         bool   `toml:"presence"`
}

type ScheduleItem struct {
	Day       string   `toml:"day"`
	Text      string   `toml:"text,omitempty"`
	Emojis    []string `toml:"emojis,omitempty"`
	Durations []string `toml:"durations,omitempty"`
}

type PresetItem struct {
	Text      string         `toml:"text"`
	Emojis    []string       `toml:"emojis,omitempty"`
	Durations []string       `toml:"durations,omitempty"`
	Schedule  []ScheduleItem `toml:"schedule,omitempty"`
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
		return Config{}, fmt.Errorf("config validation failed: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, fmt.Errorf("config validation failed: %w", err)
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

	// TODO: check emoji name

	return errs
}

func (c *Config) validateDurations(preset PresetItem, name string) []error {
	var errs []error
	var allDates []string

	if len(preset.Durations) > 0 {
		allDates = append(allDates, preset.Durations...)
	}

	for _, sch := range preset.Schedule {
		allDates = append(allDates, sch.Durations...)
	}

	for _, d := range allDates {
		if !pkg.IsValidRelativeDate(d) {
			errs = append(errs, fmt.Errorf("found invalid duration [%s] on preset [%s]", d, name))
		}
	}

	return errs
}
