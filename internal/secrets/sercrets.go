package secrets

import (
	"afk/internal/entities"
	"fmt"

	"github.com/zalando/go-keyring"
)

type Secret interface {
	SetSecret(key, value string) error
	GetSecret(key string) (string, error)
}

type secret struct {
	// Future Deps
}

func New() Secret {
	return &secret{}
}

func (s *secret) SetSecret(key, value string) error {
	err := keyring.Set(entities.KEYRING_APP_NAME, key, value)
	if err != nil {
		return fmt.Errorf("failed to save afk secret '%s': %v", key, err)
	}

	return nil
}

func (s *secret) GetSecret(key string) (string, error) {
	value, err := keyring.Get(entities.KEYRING_APP_NAME, key)
	if err != nil {
		return "", fmt.Errorf("failed to get afk secret '%s': %v", key, err)
	}
	return value, nil
}
