package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/domo84/whoop-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Defaults(t *testing.T) {
	dir := t.TempDir()
	mgr, err := config.NewManager(dir)
	require.NoError(t, err)

	cfg, err := mgr.Load()
	require.NoError(t, err)
	assert.Equal(t, "table", cfg.OutputFormat)
	assert.Equal(t, 8282, cfg.RedirectPort)
	assert.Empty(t, cfg.ClientID)
}

func TestSave_And_Load(t *testing.T) {
	dir := t.TempDir()
	mgr, err := config.NewManager(dir)
	require.NoError(t, err)

	original := &config.Config{
		OutputFormat: "json",
		ClientID:     "my-client",
		ClientSecret: "my-secret",
		RedirectPort: 9999,
	}
	require.NoError(t, mgr.Save(original))

	// Verify file exists
	cfgPath := filepath.Join(dir, "config.yaml")
	_, statErr := os.Stat(cfgPath)
	require.NoError(t, statErr)

	// Load via new manager to ensure persistence
	mgr2, err := config.NewManager(dir)
	require.NoError(t, err)
	loaded, err := mgr2.Load()
	require.NoError(t, err)

	assert.Equal(t, "json", loaded.OutputFormat)
	assert.Equal(t, "my-client", loaded.ClientID)
	assert.Equal(t, "my-secret", loaded.ClientSecret)
	assert.Equal(t, 9999, loaded.RedirectPort)
}

func TestLoad_EnvOverride(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("WHOOP_CLIENT_ID", "env-client-id")

	mgr, err := config.NewManager(dir)
	require.NoError(t, err)

	cfg, err := mgr.Load()
	require.NoError(t, err)
	assert.Equal(t, "env-client-id", cfg.ClientID)
}

func TestDir(t *testing.T) {
	dir := t.TempDir()
	mgr, err := config.NewManager(dir)
	require.NoError(t, err)
	assert.Equal(t, dir, mgr.Dir())
}
