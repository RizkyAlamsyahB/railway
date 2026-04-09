package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

// AppConfig holds application-level configuration.
type AppConfig struct {
	Name               string `mapstructure:"APP_NAME"`
	Port               int    `mapstructure:"APP_PORT"`
	Env                string `mapstructure:"APP_ENV"`
	BaseURL            string `mapstructure:"APP_BASE_URL"`
	FrontendURL        string `mapstructure:"APP_FRONTEND_URL"`
	CORSAllowedOrigins string `mapstructure:"CORS_ALLOWED_ORIGINS"`
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

// DSN returns the PostgreSQL connection string.
func (d DatabaseConfig) DSN() string {
	return d.buildURL("postgres")
}

// MigrateURL returns the database URL in the format expected by golang-migrate.
func (d DatabaseConfig) MigrateURL() string {
	return d.buildURL("postgres")
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

// OTPConfig holds OTP issuance and proof-token configuration.
type OTPConfig struct {
	CodeLength            int    `mapstructure:"OTP_CODE_LENGTH"`
	ExpiryMinutes         int    `mapstructure:"OTP_EXPIRY_MINUTES"`
	ResendCooldownSeconds int    `mapstructure:"OTP_RESEND_COOLDOWN_SECONDS"`
	MaxAttempts           int    `mapstructure:"OTP_MAX_ATTEMPTS"`
	ProofExpiryMinutes    int    `mapstructure:"OTP_PROOF_EXPIRY_MINUTES"`
	Secret                string `mapstructure:"OTP_SECRET"`
	Issuer                string `mapstructure:"OTP_ISSUER"`
}

// GoogleOAuthConfig holds Google OAuth ID token verification settings.
type GoogleOAuthConfig struct {
	ClientID     string `mapstructure:"GOOGLE_OAUTH_CLIENT_ID"`
	TokenInfoURL string `mapstructure:"GOOGLE_OAUTH_TOKENINFO_URL"`
	TimeoutSec   int    `mapstructure:"GOOGLE_OAUTH_TIMEOUT_SEC"`
}

// AdminConfig holds seed admin user configuration.
type AdminConfig struct {
	Email    string `mapstructure:"ADMIN_EMAIL"`
	Password string `mapstructure:"ADMIN_PASSWORD"`
	Name     string `mapstructure:"ADMIN_NAME"`
	Phone    string `mapstructure:"ADMIN_PHONE"`
}

// StorageConfig holds object storage configuration.
type StorageConfig struct {
	Provider         string `mapstructure:"STORAGE_PROVIDER"`
	S3Region         string `mapstructure:"STORAGE_S3_REGION"`
	S3Bucket         string `mapstructure:"STORAGE_S3_BUCKET"`
	S3AccessKey      string `mapstructure:"STORAGE_S3_ACCESS_KEY"`
	S3SecretKey      string `mapstructure:"STORAGE_S3_SECRET_KEY"`
	S3Endpoint       string `mapstructure:"STORAGE_S3_ENDPOINT"`
	S3ForcePathStyle bool   `mapstructure:"STORAGE_S3_FORCE_PATH_STYLE"`
	BaseURL          string `mapstructure:"STORAGE_BASE_URL"`
	UploadMaxSizeMB  int    `mapstructure:"STORAGE_UPLOAD_MAX_SIZE_MB"`
}

// SMTPConfig holds SMTP email sending configuration.
type SMTPConfig struct {
	Host      string `mapstructure:"SMTP_HOST"`
	Port      int    `mapstructure:"SMTP_PORT"`
	Username  string `mapstructure:"SMTP_USERNAME"`
	Password  string `mapstructure:"SMTP_PASSWORD"`
	FromEmail string `mapstructure:"SMTP_FROM_EMAIL"`
	FromName  string `mapstructure:"SMTP_FROM_NAME"`
}

// IMAPConfig holds IMAP inbox reading configuration.
type IMAPConfig struct {
	Host            string `mapstructure:"IMAP_HOST"`
	Port            int    `mapstructure:"IMAP_PORT"`
	Username        string `mapstructure:"IMAP_USERNAME"`
	Password        string `mapstructure:"IMAP_PASSWORD"`
	Enabled         bool   `mapstructure:"IMAP_ENABLED"`
	PollIntervalSec int    `mapstructure:"IMAP_POLL_INTERVAL_SEC"`
}

// XenditConfig holds Xendit payment platform configuration.
type XenditConfig struct {
	APISecretKey             string `mapstructure:"XENDIT_API_SECRET_KEY"`
	APIPublicKey             string `mapstructure:"XENDIT_API_PUBLIC_KEY"`
	BaseURL                  string `mapstructure:"XENDIT_BASE_URL"`
	WebhookVerificationToken string `mapstructure:"XENDIT_WEBHOOK_VERIFICATION_TOKEN"`
	WebhookURL               string `mapstructure:"XENDIT_WEBHOOK_URL"`
	Bypass                   bool   `mapstructure:"XENDIT_BYPASS"`
}

// RajaOngkirConfig holds RajaOngkir shipping API configuration.
type RajaOngkirConfig struct {
	APIKey  string `mapstructure:"RAJAONGKIR_API_KEY"`
	BaseURL string `mapstructure:"RAJAONGKIR_BASE_URL"`
}

// WithdrawalConfig holds vendor withdrawal fee policy configuration.
type WithdrawalConfig struct {
	FeeEstimateFixed float64 `mapstructure:"WITHDRAWAL_FEE_ESTIMATE_FIXED"`
	MinNetAmount     float64 `mapstructure:"WITHDRAWAL_MIN_NET_AMOUNT"`
}

// SensitiveDataConfig holds field-level protection settings for sensitive DB values.
type SensitiveDataConfig struct {
	EncryptionKey string `mapstructure:"SENSITIVE_DATA_ENCRYPTION_KEY"`
}

// Config is the root configuration struct containing all configuration sections.
type Config struct {
	App         AppConfig           `mapstructure:",squash"`
	Database    DatabaseConfig      `mapstructure:",squash"`
	JWT         JWTConfig           `mapstructure:",squash"`
	OTP         OTPConfig           `mapstructure:",squash"`
	Admin       AdminConfig         `mapstructure:",squash"`
	Storage     StorageConfig       `mapstructure:",squash"`
	SMTP        SMTPConfig          `mapstructure:",squash"`
	IMAP        IMAPConfig          `mapstructure:",squash"`
	Xendit      XenditConfig        `mapstructure:",squash"`
	Withdrawal  WithdrawalConfig    `mapstructure:",squash"`
	Sensitive   SensitiveDataConfig `mapstructure:",squash"`
	RajaOngkir  RajaOngkirConfig    `mapstructure:",squash"`
	GoogleOAuth GoogleOAuthConfig   `mapstructure:",squash"`
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
	viper.SetDefault("APP_BASE_URL", "http://localhost:8080")
	viper.SetDefault("APP_FRONTEND_URL", "http://localhost:3000")
	viper.SetDefault("CORS_ALLOWED_ORIGINS", "*")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "haji_umroh_store")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("DB_MAX_IDLE_CONNS", 10)
	viper.SetDefault("DB_MAX_OPEN_CONNS", 100)
	viper.SetDefault("JWT_SECRET", "")
	viper.SetDefault("JWT_EXPIRY_HOURS", 24)
	viper.SetDefault("JWT_ISSUER", "haji-umroh-store-be")
	viper.SetDefault("OTP_CODE_LENGTH", 6)
	viper.SetDefault("OTP_EXPIRY_MINUTES", 5)
	viper.SetDefault("OTP_RESEND_COOLDOWN_SECONDS", 60)
	viper.SetDefault("OTP_MAX_ATTEMPTS", 5)
	viper.SetDefault("OTP_PROOF_EXPIRY_MINUTES", 10)
	viper.SetDefault("OTP_SECRET", "")
	viper.SetDefault("OTP_ISSUER", "haji-umroh-store-be-otp")
	viper.SetDefault("GOOGLE_OAUTH_CLIENT_ID", "")
	viper.SetDefault("GOOGLE_OAUTH_TOKENINFO_URL", "https://oauth2.googleapis.com/tokeninfo")
	viper.SetDefault("GOOGLE_OAUTH_TIMEOUT_SEC", 10)
	viper.SetDefault("ADMIN_EMAIL", "")
	viper.SetDefault("ADMIN_PASSWORD", "")
	viper.SetDefault("ADMIN_NAME", "")
	viper.SetDefault("ADMIN_PHONE", "")
	viper.SetDefault("STORAGE_PROVIDER", "s3")
	viper.SetDefault("STORAGE_S3_REGION", "ap-southeast-1")
	viper.SetDefault("STORAGE_S3_BUCKET", "")
	viper.SetDefault("STORAGE_S3_ACCESS_KEY", "")
	viper.SetDefault("STORAGE_S3_SECRET_KEY", "")
	viper.SetDefault("STORAGE_S3_ENDPOINT", "")
	viper.SetDefault("STORAGE_S3_FORCE_PATH_STYLE", false)
	viper.SetDefault("STORAGE_BASE_URL", "")
	viper.SetDefault("STORAGE_UPLOAD_MAX_SIZE_MB", 10)
	viper.SetDefault("SMTP_HOST", "smtp.gmail.com")
	viper.SetDefault("SMTP_PORT", 587)
	viper.SetDefault("SMTP_USERNAME", "")
	viper.SetDefault("SMTP_PASSWORD", "")
	viper.SetDefault("SMTP_FROM_EMAIL", "")
	viper.SetDefault("SMTP_FROM_NAME", "")
	viper.SetDefault("IMAP_HOST", "imap.gmail.com")
	viper.SetDefault("IMAP_PORT", 993)
	viper.SetDefault("IMAP_USERNAME", "")
	viper.SetDefault("IMAP_PASSWORD", "")
	viper.SetDefault("IMAP_ENABLED", false)
	viper.SetDefault("IMAP_POLL_INTERVAL_SEC", 300)
	viper.SetDefault("XENDIT_API_SECRET_KEY", "")
	viper.SetDefault("XENDIT_API_PUBLIC_KEY", "")
	viper.SetDefault("XENDIT_BASE_URL", "https://api.xendit.co")
	viper.SetDefault("XENDIT_WEBHOOK_VERIFICATION_TOKEN", "")
	viper.SetDefault("XENDIT_WEBHOOK_URL", "")
	viper.SetDefault("XENDIT_BYPASS", false)
	viper.SetDefault("WITHDRAWAL_FEE_ESTIMATE_FIXED", 0)
	viper.SetDefault("WITHDRAWAL_MIN_NET_AMOUNT", 10000)
	viper.SetDefault("SENSITIVE_DATA_ENCRYPTION_KEY", "")
	viper.SetDefault("RAJAONGKIR_API_KEY", "")
	viper.SetDefault("RAJAONGKIR_BASE_URL", "https://rajaongkir.komerce.id/api/v1")

	// Read .env file (ignore error if file doesn't exist)
	_ = viper.ReadInConfig()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
