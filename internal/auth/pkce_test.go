package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPKCEChallenge(t *testing.T) {
	ch, err := NewPKCEChallenge()
	require.NoError(t, err)

	assert.NotEmpty(t, ch.Verifier)
	assert.NotEmpty(t, ch.Challenge)
	assert.Equal(t, "S256", ch.Method)

	// Verify the S256 challenge is correctly derived from the verifier.
	h := sha256.Sum256([]byte(ch.Verifier))
	expected := base64.RawURLEncoding.EncodeToString(h[:])
	assert.Equal(t, expected, ch.Challenge)
}

func TestNewPKCEChallenge_Uniqueness(t *testing.T) {
	ch1, err := NewPKCEChallenge()
	require.NoError(t, err)
	ch2, err := NewPKCEChallenge()
	require.NoError(t, err)

	assert.NotEqual(t, ch1.Verifier, ch2.Verifier, "verifiers must be unique across calls")
	assert.NotEqual(t, ch1.Challenge, ch2.Challenge, "challenges must be unique across calls")
}

func TestGenerateState(t *testing.T) {
	s1, err := generateState()
	require.NoError(t, err)
	s2, err := generateState()
	require.NoError(t, err)

	assert.NotEmpty(t, s1)
	assert.NotEmpty(t, s2)
	assert.NotEqual(t, s1, s2, "states must be unique")
}
