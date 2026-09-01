package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Binance  BinanceConfig  `mapstructure:"binance"`
	Ingest   IngestConfig   `mapstructure:"ingest"`
}

type AppConfig struct {
	Env             string        `mapstructure:"env"`
	Port            int           `mapstructure:"port"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"db"`
	SSLMode  string `mapstructure:"sslmode"`

	MaxConns int32 `mapstructure:"max_conns"`
	MinConns int32 `mapstructure:"min_conns"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	DB   int    `mapstructure:"db"`
}

type BinanceConfig struct {
	// Empty falls back to the client's own default endpoint.
	WSURL string `mapstructure:"ws_url"`
}

// The process role is chosen by which binary runs, not by config.
type IngestConfig struct {
	Symbols    []string      `mapstructure:"symbols"`
	MaxSilence time.Duration `mapstructure:"max_silence"`
}

// Load reads configuration in order of precedence: environment > .env > defaults.
func Load() (*Config, error) {
	if err := loadDotEnv(); err != nil {
		return nil, err
	}

	v := viper.New()
	setDefaults(v)

	// Maps "postgres.max_conns" to POSTGRES_MAX_CONNS.
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// loadDotEnv loads .env into the process environment without overriding
// existing variables. The binary starts in backend/, the file lives one level up.
func loadDotEnv() error {
	for _, dir := range []string{".", ".."} {
		path := filepath.Join(dir, ".env")
		if _, err := os.Stat(path); err != nil {
			continue
		}
		if err := gotenv.Load(path); err != nil {
			return fmt.Errorf("load %s: %w", path, err)
		}
		return nil
	}
	return nil // a missing .env is valid: production uses the environment
}

// setDefaults registers every known key. Unmarshal only fills keys Viper
// already knows, so a key without a default stays empty even if its
// environment variable is set.
func setDefaults(v *viper.Viper) {
	v.SetDefault("app.env", "dev")
	v.SetDefault("app.port", 8080)
	v.SetDefault("app.shutdown_timeout", "10s")

	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", 5432)
	v.SetDefault("postgres.user", "")     // required
	v.SetDefault("postgres.password", "") // required
	v.SetDefault("postgres.db", "tickstorm")
	v.SetDefault("postgres.sslmode", "disable")
	v.SetDefault("postgres.max_conns", 50)
	v.SetDefault("postgres.min_conns", 5)

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)

	v.SetDefault("binance.ws_url", "")

	// Viper splits a comma-separated env var into a slice.
	v.SetDefault("ingest.symbols", []string{"BTCUSDT", "ETHUSDT"})
	v.SetDefault("ingest.max_silence", "30s")
}

func (c *Config) validate() error {
	if c.Postgres.User == "" {
		return errors.New("POSTGRES_USER is required")
	}
	if c.Postgres.Password == "" {
		return errors.New("POSTGRES_PASSWORD is required")
	}
	if c.App.Port < 1 || c.App.Port > 65535 {
		return fmt.Errorf("invalid APP_PORT: %d", c.App.Port)
	}
	if len(c.Ingest.Symbols) == 0 {
		return errors.New("INGEST_SYMBOLS is required")
	}
	if c.Postgres.MinConns > c.Postgres.MaxConns {
		return fmt.Errorf(
			"POSTGRES_MIN_CONNS (%d) cannot exceed POSTGRES_MAX_CONNS (%d)",
			c.Postgres.MinConns,
			c.Postgres.MaxConns,
		)
	}
	return nil
}

// DSN returns the Postgres connection string for pgxpool.
func (c *PostgresConfig) DSN() string {
	u := &url.URL{
		Scheme:   "postgresql",
		User:     url.UserPassword(c.User, c.Password),
		Host:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Path:     c.Database,
		RawQuery: url.Values{"sslmode": {c.SSLMode}}.Encode(),
	}
	return u.String()
}

func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *AppConfig) IsProduction() bool {
	return c.Env == "prod"
}
