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
	PickDuration string `toml:"pickDuration"`
}

type PresetItem struct {
	Text      string   `toml:"text"`
	Emojis    []string `toml:"emojis,omitempty"`
	Durations []string `toml:"durations"`
}

func New() (Config, error) {
	cfgFile, err := pkg.GetConfigPath(entities.AFK_CONFIG_FILENAME)
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		panic(fmt.Errorf("failed to read file: %w", err))
	}

	var cfg Config
	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal TOML: %w", err))
	}

	// fmt.Printf("Pick Duration Strategy: %s\n", cfg.User.PickDuration)

	// fmt.Println("\nAvailable Presets:")
	// for name, preset := range cfg.Presets {
	// 	fmt.Printf("- [%s]: %s (Durations: %v)\n", name, preset.Text, preset.Durations)
	// }

	return cfg, nil
}
