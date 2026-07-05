package config

import (
	"afk/internal/entities"
	"afk/pkg"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

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
	// TODO: validate emoji strings
	// TODO: validate duration strings

	const maxLen = entities.SLACK_MAX_TEXT_LENGTH

	for name, preset := range c.Presets {
		// Check default preset text length
		if len(preset.Text) > maxLen {
			return fmt.Errorf("preset [%s] text exceeds max length of %d (got %d chars)", name, maxLen, len(preset.Text))
		}

		// Check each day-specific schedule text length as well
		for _, schedule := range preset.Schedule {
			if len(schedule.Text) > maxLen {
				return fmt.Errorf("preset [%s] schedule for %s text exceeds max length of %d (got %d chars)", name, schedule.Day, maxLen, len(schedule.Text))
			}
		}
	}

	return nil
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
