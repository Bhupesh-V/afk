package secrets

type Secret interface {
	// GetTokens interacts with user's keyring and orchestrates the token getter logic
	GetTokens() (map[string]string, error)
	// UpdateTokens appends any new token and updates any existing tokens to user's keyring
	UpdateTokens(map[string]string) error
}

type secret struct {
}

func New() Secret {
	return &secret{}
}

func (s *secret) GetTokens() (map[string]string, error) {
	return nil, nil
}
func (s *secret) UpdateTokens(map[string]string) error {
	return nil
}
