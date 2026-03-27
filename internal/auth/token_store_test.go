package auth

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestFileTokenStore_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store := NewFileTokenStore(dir)

	token := &oauth2.Token{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(1 * time.Hour).Truncate(time.Second),
	}

	require.NoError(t, store.Save(token))
	assert.True(t, store.Exists())

	// Verify file permissions are 0600.
	info, err := os.Stat(filepath.Join(dir, "credentials.json"))
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())

	loaded, err := store.Load()
	require.NoError(t, err)
	assert.Equal(t, token.AccessToken, loaded.AccessToken)
	assert.Equal(t, token.RefreshToken, loaded.RefreshToken)
	assert.Equal(t, token.TokenType, loaded.TokenType)
}

func TestFileTokenStore_Exists_False_WhenMissing(t *testing.T) {
	dir := t.TempDir()
	store := NewFileTokenStore(dir)
	assert.False(t, store.Exists())
}

func TestFileTokenStore_Load_NotFound(t *testing.T) {
	dir := t.TempDir()
	store := NewFileTokenStore(dir)
	_, err := store.Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not authenticated")
}

func TestFileTokenStore_Delete(t *testing.T) {
	dir := t.TempDir()
	store := NewFileTokenStore(dir)
	token := &oauth2.Token{AccessToken: "x"}

	require.NoError(t, store.Save(token))
	require.True(t, store.Exists())

	require.NoError(t, store.Delete())
	assert.False(t, store.Exists())
}

func TestFileTokenStore_Delete_NonExistent(t *testing.T) {
	dir := t.TempDir()
	store := NewFileTokenStore(dir)
	// Deleting when file doesn't exist should not error.
	assert.NoError(t, store.Delete())
}
