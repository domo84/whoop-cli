package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

// PKCEChallenge holds a PKCE verifier and its corresponding S256 challenge.
type PKCEChallenge struct {
	Verifier  string
	Challenge string
	Method    string
}

// NewPKCEChallenge generates a cryptographically random verifier and derives
// the S256 challenge from it.
func NewPKCEChallenge() (*PKCEChallenge, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	verifier := base64.RawURLEncoding.EncodeToString(b)

	h := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(h[:])

	return &PKCEChallenge{
		Verifier:  verifier,
		Challenge: challenge,
		Method:    "S256",
	}, nil
}

// generateState produces a random string suitable for use as an OAuth2 state parameter.
func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
