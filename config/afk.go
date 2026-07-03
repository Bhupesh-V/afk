package config

import "fmt"

// Status represents the configuration for a single presence status.
type Status struct {
	Name     string `yaml:"name"`
	Emoji    string `yaml:"emoji"`
	Text     string `yaml:"text"`
	Duration string `yaml:"duration"` // Read as string to handle formats like "60m"
}

// StatusConfig maps to the root key of your YAML file.
type StatusConfig struct {
	Statuses []Status `yaml:"statuses"`
}

func (s *StatusConfig) Validate() error {
	for _, status := range s.Statuses {
		if status.Name == "" {
			return fmt.Errorf("status name cannot be empty")
		}
		if status.Emoji == "" {
			return fmt.Errorf("emoji for status %s cannot be empty", status.Name)
		}
		if status.Text == "" {
			return fmt.Errorf("text for status %s cannot be empty", status.Name)
		}
		if status.Duration == "" {
			return fmt.Errorf("duration for status %s cannot be empty", status.Name)
		}
	}
	return nil
}
