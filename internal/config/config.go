package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// parseEnvList splits a comma-separated env var into a string slice, trimming spaces
func parseEnvList(key string) []string {
	val := viper.GetString(key)
	if val == "" {
		return nil
	}
	var out []string
	for _, v := range strings.Split(val, ",") {
		v = strings.TrimSpace(v)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}


type Config struct {
	App      AppConfig
	Log      LogConfig
	CORS     CORSConfig
	API      APIConfig
	Network  NetworkConfig
	Services ServicesConfig
	Auth     AuthConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port string
	Host string
}

type LogConfig struct {
	Level  string
	Format string
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

type APIConfig struct {
	Prefix  string
	Timeout time.Duration
}

type NetworkConfig struct {
	DefaultInterface string
	NetplanPath      string
}

type ServicesConfig struct {
	LogstashPath string
	SuricataPath string
	RsyslogPath  string
	WazuhPath    string
	RPortPath    string
}

type AuthConfig struct {
	JWTSecret     string
	JWTExpiration time.Duration
}

// Load loads the application configuration
func Load() (*Config, error) {
	// Load .env file
	_ = godotenv.Load()

	// Set defaults
	viper.SetDefault("APP_NAME", "NTA Backend")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_HOST", "0.0.0.0")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "json")
	viper.SetDefault("API_PREFIX", "/api/v1")
	viper.SetDefault("API_TIMEOUT", "30s")
	viper.SetDefault("DEFAULT_NETWORK_INTERFACE", "enp3s0")
	viper.SetDefault("NETPLAN_CONFIG_PATH", "")
       netplanPath := viper.GetString("NETPLAN_CONFIG_PATH")
       if netplanPath == "" {
	       // Scan /etc/netplan for YAML files
	       files, err := os.ReadDir("/etc/netplan")
	       if err == nil {
		       for _, f := range files {
			       if !f.IsDir() && (strings.HasSuffix(f.Name(), ".yaml") || strings.HasSuffix(f.Name(), ".yml")) {
				       netplanPath = "/etc/netplan/" + f.Name()
				       break
			       }
		       }
	       }
       }

	// Bind environment variables
	viper.AutomaticEnv()

	timeout, err := time.ParseDuration(viper.GetString("API_TIMEOUT"))
	if err != nil {
		timeout = 30 * time.Second
	}

	       // JWT_EXPIRATION is not used here

	// Parse comma-separated CORS env vars into slices
	allowedOrigins := parseEnvList("CORS_ALLOWED_ORIGINS")
	allowedMethods := parseEnvList("CORS_ALLOWED_METHODS")
	allowedHeaders := parseEnvList("CORS_ALLOWED_HEADERS")

	       config := &Config{
		       App: AppConfig{
			       Name: viper.GetString("APP_NAME"),
			       Env:  viper.GetString("APP_ENV"),
			       Port: viper.GetString("APP_PORT"),
			       Host: viper.GetString("APP_HOST"),
		       },
		       Log: LogConfig{
			       Level:  viper.GetString("LOG_LEVEL"),
			       Format: viper.GetString("LOG_FORMAT"),
		       },
		       CORS: CORSConfig{
			       AllowedOrigins: allowedOrigins,
			       AllowedMethods: allowedMethods,
			       AllowedHeaders: allowedHeaders,
		       },
		       API: APIConfig{
			       Prefix:  viper.GetString("API_PREFIX"),
			       Timeout: timeout,
		       },
		       Network: NetworkConfig{
			       DefaultInterface: viper.GetString("DEFAULT_NETWORK_INTERFACE"),
			       NetplanPath:      netplanPath,
		       },
		       Services: ServicesConfig{
			       LogstashPath: viper.GetString("LOGSTASH_CONFIG_PATH"),
		       },
	       }

		if err := config.Validate(); err != nil {
			return nil, fmt.Errorf("invalid configuration: %w", err)
		}

		return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.App.Port == "" {
		return fmt.Errorf("APP_PORT is required")
	}

	if c.App.Name == "" {
		return fmt.Errorf("APP_NAME is required")
	}

	return nil
}

// IsDevelopment returns true if the app is running in development mode
func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development" || c.App.Env == "dev"
}

// IsProduction returns true if the app is running in production mode
func (c *Config) IsProduction() bool {
	return c.App.Env == "production" || c.App.Env == "prod"
}
