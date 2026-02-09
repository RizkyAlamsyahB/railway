package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

// AppConfig holds application-level configuration.
type AppConfig struct {
	Name string `mapstructure:"APP_NAME"`
	Port int    `mapstructure:"APP_PORT"`
	Env  string `mapstructure:"APP_ENV"`
}

// DatabaseConfig holds database connection configuration.
type DatabaseConfig struct {
	Host         string `mapstructure:"DB_HOST"`
	Port         int    `mapstructure:"DB_PORT"`
	User         string `mapstructure:"DB_USER"`
	Password     string `mapstructure:"DB_PASSWORD"`
	Name         string `mapstructure:"DB_NAME"`
	SSLMode      string `mapstructure:"DB_SSLMODE"`
	MaxIdleConns int    `mapstructure:"DB_MAX_IDLE_CONNS"`
	MaxOpenConns int    `mapstructure:"DB_MAX_OPEN_CONNS"`
}

// DSN returns the PostgreSQL-compatible connection string for CockroachDB.
func (d DatabaseConfig) DSN() string {
	return d.buildURL("postgres")
}

// MigrateURL returns the database URL in the format expected by golang-migrate.
func (d DatabaseConfig) MigrateURL() string {
	return d.buildURL("cockroachdb")
}

func (d DatabaseConfig) buildURL(scheme string) string {
	u := &url.URL{
		Scheme: scheme,
		Host:   fmt.Sprintf("%s:%d", d.Host, d.Port),
		Path:   "/" + d.Name,
	}

	if d.Password == "" {
		u.User = url.User(d.User)
	} else {
		u.User = url.UserPassword(d.User, d.Password)
	}

	q := u.Query()
	q.Set("sslmode", d.SSLMode)
	u.RawQuery = q.Encode()

	return u.String()
}

// JWTConfig holds JSON Web Token configuration.
type JWTConfig struct {
	Secret      string `mapstructure:"JWT_SECRET"`
	ExpiryHours int    `mapstructure:"JWT_EXPIRY_HOURS"`
	Issuer      string `mapstructure:"JWT_ISSUER"`
}

// AdminConfig holds seed admin user configuration.
type AdminConfig struct {
	Email    string `mapstructure:"ADMIN_EMAIL"`
	Password string `mapstructure:"ADMIN_PASSWORD"`
	Name     string `mapstructure:"ADMIN_NAME"`
	Phone    string `mapstructure:"ADMIN_PHONE"`
}

// Config is the root configuration struct containing all configuration sections.
type Config struct {
	App      AppConfig      `mapstructure:",squash"`
	Database DatabaseConfig `mapstructure:",squash"`
	JWT      JWTConfig      `mapstructure:",squash"`
	Admin    AdminConfig    `mapstructure:",squash"`
}

// Load reads configuration from the .env file and environment variables.
// Environment variables take precedence over .env file values.
func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Defaults
	viper.SetDefault("APP_NAME", "haji-umroh-store-be")
	viper.SetDefault("APP_PORT", 8080)
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 26257)
	viper.SetDefault("DB_USER", "root")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_NAME", "haji_umroh_store")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("DB_MAX_IDLE_CONNS", 10)
	viper.SetDefault("DB_MAX_OPEN_CONNS", 100)
	viper.SetDefault("JWT_SECRET", "")
	viper.SetDefault("JWT_EXPIRY_HOURS", 24)
	viper.SetDefault("JWT_ISSUER", "haji-umroh-store-be")
	viper.SetDefault("ADMIN_EMAIL", "")
	viper.SetDefault("ADMIN_PASSWORD", "")
	viper.SetDefault("ADMIN_NAME", "")
	viper.SetDefault("ADMIN_PHONE", "")

	// Read .env file (ignore error if file doesn't exist)
	_ = viper.ReadInConfig()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
