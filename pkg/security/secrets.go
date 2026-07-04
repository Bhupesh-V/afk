package security

import (
	"afk/internal/entities"
	"log"

	"github.com/zalando/go-keyring"
)

func SetSecret(key, value string) error {
	err := keyring.Set(entities.KEYRING_APP_NAME, key, value)
	if err != nil {
		log.Fatalf("Failed to save afk secret '%s' :%v", key, err)
	}

	return nil
}

func GetSecret(key string) (string, error) {
	value, err := keyring.Get(entities.KEYRING_APP_NAME, key)
	if err != nil {
		log.Fatalf("Failed to get afk secret '%s' :%v", key, err)
	}
	return value, nil
}
