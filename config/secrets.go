package config

import "fmt"

// Account represents an individual team configuration.
type Account struct {
	Team   string `yaml:"team"`
	Secret string `yaml:"secret"`
}

// Config maps to the root key of the YAML file.
type Config struct {
	Accounts []Account `yaml:"accounts"`
}

func (c *Config) Validate() error {
	for _, account := range c.Accounts {
		if account.Team == "" {
			return fmt.Errorf("team name cannot be empty")
		}
		if account.Secret == "" {
			return fmt.Errorf("secret for team %s cannot be empty", account.Team)
		}
	}
	return nil
}
