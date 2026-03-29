package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	defaultConfigDir  = ".config/whoop-cli"
	defaultConfigFile = "config"
	defaultPort       = 8282
	defaultOutput     = "table"
)

// Config holds all application settings.
type Config struct {
	OutputFormat string `mapstructure:"output_format" yaml:"output_format"`
	ClientID     string `mapstructure:"client_id"     yaml:"client_id"`
	ClientSecret string `mapstructure:"client_secret" yaml:"client_secret"`
	RedirectPort int    `mapstructure:"redirect_port" yaml:"redirect_port"`
}

// Manager handles reading and writing the application config file.
type Manager interface {
	Load() (*Config, error)
	Save(cfg *Config) error
	Dir() string
}

type viperManager struct {
	v      *viper.Viper
	dir    string
	cfgPath string
}

// NewManager creates a Manager rooted at configDir.
// If configDir is empty, it defaults to ~/.whoop-cli/.
func NewManager(configDir string) (Manager, error) {
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("could not determine home directory: %w", err)
		}
		configDir = filepath.Join(home, defaultConfigDir)
	}

	v := viper.New()
	v.SetDefault("output_format", defaultOutput)
	v.SetDefault("redirect_port", defaultPort)
	v.SetConfigName(defaultConfigFile)
	v.SetConfigType("yaml")
	v.AddConfigPath(configDir)
	v.SetEnvPrefix("WHOOP")
	v.AutomaticEnv()
	// Explicit bindings ensure env vars are picked up during Unmarshal.
	v.BindEnv("client_id", "WHOOP_CLIENT_ID")         //nolint:errcheck
	v.BindEnv("client_secret", "WHOOP_CLIENT_SECRET") //nolint:errcheck
	v.BindEnv("output_format", "WHOOP_OUTPUT")        //nolint:errcheck
	v.BindEnv("redirect_port", "WHOOP_REDIRECT_PORT") //nolint:errcheck

	return &viperManager{
		v:       v,
		dir:     configDir,
		cfgPath: filepath.Join(configDir, defaultConfigFile+".yaml"),
	}, nil
}

func (m *viperManager) Load() (*Config, error) {
	if err := m.v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config: %w", err)
		}
		// Config file not found is acceptable — use defaults and env vars.
	}
	var cfg Config
	if err := m.v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshalling config: %w", err)
	}
	return &cfg, nil
}

func (m *viperManager) Save(cfg *Config) error {
	if err := os.MkdirAll(m.dir, 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}
	m.v.Set("output_format", cfg.OutputFormat)
	m.v.Set("redirect_port", cfg.RedirectPort)
	// Never persist client_id / client_secret to disk if they came from env.
	if cfg.ClientID != "" {
		m.v.Set("client_id", cfg.ClientID)
	}
	if cfg.ClientSecret != "" {
		m.v.Set("client_secret", cfg.ClientSecret)
	}
	return m.v.WriteConfigAs(m.cfgPath)
}

func (m *viperManager) Dir() string {
	return m.dir
}
