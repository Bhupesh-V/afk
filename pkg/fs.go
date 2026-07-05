package pkg

import (
	"afk/internal/entities"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// GetConfigPath returns the absolute path for the config file inside the "afk" directory.
// It overrides macOS to use ~/.config, while leaving Linux (~/.config) and Windows (%APPDATA%) intact.
// It creates the directory structure and the file if they do not exist.
func GetConfigPath(filename string) (string, error) {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	// Enforce strict XDG naming on macOS, fuck "Application Support"
	if runtime.GOOS == "darwin" {
		// Respect $XDG_CONFIG_HOME if explicitly set by the mac user
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			baseDir = xdg
		} else if home, err := os.UserHomeDir(); err == nil {
			baseDir = filepath.Join(home, ".config")
		}
	}

	configDir := filepath.Join(baseDir, entities.AFK_CONFIG_DIR_NAME)

	// Ensure the parent directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	fullPath := filepath.Join(configDir, filename)

	// Ensure the actual file exists without truncating/clearing it if it already has data
	file, err := os.OpenFile(fullPath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create config file: %w", err)
	}
	file.Close()

	return fullPath, nil
}
