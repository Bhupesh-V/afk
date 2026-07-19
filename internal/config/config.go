package config

import (
	"afk/internal/entities"
	"afk/pkg"
	"embed"
	"errors"
	"fmt"
	"os"
	"strings"

	slackEntities "afk/internal/clients/slack/entities"

	"github.com/mrosales/emoji-go"
	"github.com/pelletier/go-toml/v2"
)

//go:embed sample.toml
var DefaultConfig embed.FS

type Config struct {
	emojiMap map[string]emoji.Info
	User     UserSection           `toml:"user"`
	Presets  map[string]PresetItem `toml:"presets"`
}

type UserSection struct {
	Provider string `toml:"provider" comment:"The messaging provider you use (e.g. slack)"`
}

type ScheduleItem struct {
	Day    string   `toml:"day"`
	Text   string   `toml:"text,omitempty"`
	Emojis []string `toml:"emojis,omitempty" multiline:"true"`
}

type PresetItem struct {
	Text          string         `toml:"text"`
	Emojis        []string       `toml:"emojis,omitempty" multiline:"true"`
	Duration      string         `toml:"duration,omitempty"`
	Presence      string         `toml:"presence,omitempty"`
	Notifications string         `toml:"notifications,omitempty"`
	Schedule      []ScheduleItem `toml:"schedule,omitempty" inline:"true"`
}

// Validate ensures that no text field in defaults or schedules exceeds 100 characters.
func (c *Config) Validate() error {
	var errs []error
	var provider = strings.ToLower(strings.TrimSpace(c.User.Provider))

	switch provider {
	case entities.PROVIDER_SLACK:
	case "":
		errs = append(errs, fmt.Errorf("missing 'provider' in config"))
	default:
		errs = append(errs, fmt.Errorf("unsupported provider [%s] in config", c.User.Provider))
	}

	c.User.Provider = provider

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

	emojiLookup := make(map[string]emoji.Info)
	for _, de := range emoji.All {
		// Map by the main name
		emojiLookup[de.Name] = de

		// Map by all alternate names as well
		for _, name := range de.AlternateNames {
			emojiLookup[name] = de
		}
	}

	cfg.emojiMap = emojiLookup

	if cfg.IsEmpty() {
		data, err = DefaultConfig.ReadFile("sample.toml")
		if err != nil {
			return Config{}, fmt.Errorf("failed to read embedded sample.toml: %w", err)
		}

		// Write the embedded template back to disk for the user's future runs
		err = os.WriteFile(cfgFile, data, 0600)
		if err != nil {
			return Config{}, fmt.Errorf("failed to initialize empty config file: %w", err)
		}

		// Parse the fallback config into our struct
		err = toml.Unmarshal(data, &cfg)
		if err != nil {
			return Config{}, fmt.Errorf("failed to unmarshal default config: %w", err)
		}
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, fmt.Errorf("afk config validation failed:\n%w", err)
	}

	return cfg, nil
}

// Write marshals the current Config back to the config file with proper formatting.
func (c *Config) Write() error {
	cfgFile, err := pkg.GetConfigPath(entities.AFK_CONFIG_FILENAME)
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	file, err := os.OpenFile(cfgFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	encoder.SetArraysMultiline(true) // Multi-line arrays

	err = encoder.Encode(c)
	if err != nil {
		return fmt.Errorf("failed to write and format config: %w", err)
	}

	return nil
}

// IsEmpty returns true if the configuration has no provider and no presets.
func (c *Config) IsEmpty() bool {
	return c.User.Provider == "" && len(c.Presets) == 0
}

func (c *Config) validateText(preset PresetItem, name string) []error {
	const maxLen = slackEntities.SLACK_MAX_TEXT_LENGTH
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

	for _, e := range emojis {
		if _, found := c.emojiMap[e]; !found {
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
