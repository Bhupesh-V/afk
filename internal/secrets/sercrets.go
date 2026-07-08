package secrets

import (
	"afk/internal/entities"
	keyring "afk/pkg/security"
	"fmt"
	"strings"
)

type Secret interface {
	// GetTokens interacts with user's keyring and orchestrates the token getter logic
	GetTokens(provider string) (map[string]string, error)
	// UpdateTokens appends any new token and updates any existing tokens to user's keyring
	UpdateTokens(tokens map[string]entities.SlackMetaData, provider string) error
}

type secret struct {
	// Future Deps
}

func New() Secret {
	return &secret{}
}

func (s *secret) GetTokens(provider string) (map[string]string, error) {
	tokens := make(map[string]string)

	rawTeamIds, err := keyring.GetSecret(entities.KEYRING_SLACK_TEAM_IDS)
	if err != nil {
		return nil, err
	}

	for _, id := range strings.Split(rawTeamIds, ",") {
		tokenKey := fmt.Sprintf("%s_%s", entities.KEYRING_TEAM_SECRET_PREFIX, id)
		token, err := keyring.GetSecret(tokenKey)
		if err != nil {
			return nil, err
		}
		// rawTokens[0] is team namere
		// rawTokens[1] is token
		rawTokens := strings.Split(token, entities.KEYRING_SECRET_SEP)
		if len(rawTokens) != 2 {
			return nil, fmt.Errorf("inconsisten token format on system keyring")
		}
		tokens[rawTokens[0]] = rawTokens[1]
	}

	return nil, nil
}

func (s *secret) UpdateTokens(tokens map[string]entities.SlackMetaData, provider string) error {
	rawTeamIds, err := keyring.GetSecret(entities.KEYRING_SLACK_TEAM_IDS)
	if err != nil {
		return err
	}
	teamIds := strings.Split(rawTeamIds, ",")

	for teamId, meta := range tokens {
		var found bool
		for _, id := range teamIds {
			if id == teamId {
				break
			}
		}
		tokenKey := fmt.Sprintf("%s_%s", entities.KEYRING_TEAM_SECRET_PREFIX, teamId)

		if !found {
			// save the team id in the team list secret
			teamIds = append(teamIds, teamId)
			err = keyring.SetSecret(entities.KEYRING_SLACK_TEAM_IDS, strings.Join(teamIds, ","))
			if err != nil {
				return err
			}
		}

		// overwrite/create the secret
		err = keyring.SetSecret(tokenKey, fmt.Sprintf("%s|%s", meta.TeamName, meta.WorkspaceToken))
		if err != nil {
			return err
		}
	}

	return nil
}
