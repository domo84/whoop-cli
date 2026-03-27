package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
)

const credentialsFile = "credentials.json"

// TokenStore persists and retrieves OAuth2 tokens.
type TokenStore interface {
	Save(token *oauth2.Token) error
	Load() (*oauth2.Token, error)
	Delete() error
	Exists() bool
}

type fileTokenStore struct {
	path string
}

// NewFileTokenStore creates a TokenStore that reads/writes credentials.json
// inside configDir.
func NewFileTokenStore(configDir string) TokenStore {
	return &fileTokenStore{
		path: filepath.Join(configDir, credentialsFile),
	}
}

func (s *fileTokenStore) Save(token *oauth2.Token) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return fmt.Errorf("creating credentials directory: %w", err)
	}
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling token: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0600); err != nil {
		return fmt.Errorf("writing credentials: %w", err)
	}
	return nil
}

func (s *fileTokenStore) Load() (*oauth2.Token, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("not authenticated: run 'whoop auth login' first")
		}
		return nil, fmt.Errorf("reading credentials: %w", err)
	}
	var token oauth2.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("parsing credentials: %w", err)
	}
	return &token, nil
}

func (s *fileTokenStore) Delete() error {
	if err := os.Remove(s.path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("deleting credentials: %w", err)
	}
	return nil
}

func (s *fileTokenStore) Exists() bool {
	_, err := os.Stat(s.path)
	return err == nil
}
