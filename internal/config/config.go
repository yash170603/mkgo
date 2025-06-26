package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseConfig DatabaseConfig `mapstructure:"database"`
	JwtConfig      JwtConfig      `mapstructure:"auth"`
	APIConfig      APIConfig      `mapstructure:"api"`
}

type APIConfig struct {
	Port int `mapstructure:"port"`
}

type JwtConfig struct {
	JWTSecret string `mapstructure:"secret_key"`
	ExpiresIn int    `mapstructure:"expires_in"`
	JWTIssuer string `mapstructure:"issuer"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
}

func New() Config {
	v := viper.New()
	v.SetConfigType("yaml")
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found or error loading it", "error", err.Error())
		// Don't exit - environment variables might be set through other means
	}

	// Set default config path if CONFIG_PATH is not set
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}
	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		slog.Error("Failed to read config file", slog.String("error", err.Error()), slog.String("path", configPath))
		os.Exit(1)
	}

	cfg := Config{}
	if err := v.Unmarshal(&cfg); err != nil {
		slog.Error("Failed to unmarshal config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	cfg.PopulateEnv()

	slog.Info(" Config Initialized")
	//print the whole config
	slog.Info("Config", "config", cfg)
	return cfg
}

func (c *Config) PopulateEnv() {
	// API
	if env := os.Getenv("API_PORT"); env != "" {
		if port, err := strconv.Atoi(env); err == nil {
			c.APIConfig.Port = port
		} else {
			slog.Error("Invalid API_PORT", "value", env, "error", err.Error())
		}
	}

	if env := os.Getenv("JWT_SECRET"); env != "" {
		c.JwtConfig.JWTSecret = env
	}
	if env := os.Getenv("JWT_ISSUER"); env != "" {
		c.JwtConfig.JWTIssuer = env
	}
	if env := os.Getenv("JWT_EXPIRES_IN"); env != "" {
		if expiresIn, err := strconv.Atoi(env); err == nil {
			c.JwtConfig.ExpiresIn = expiresIn
		} else {
			slog.Error("Invalid JWT_EXPIRES_IN", "value", env, "error", err.Error())
		}
	}

	if env := os.Getenv("DB_HOST"); env != "" {
		c.DatabaseConfig.Host = env
	}
	if env := os.Getenv("DB_USER"); env != "" {
		c.DatabaseConfig.User = env
	}
	if env := os.Getenv("DB_PASSWORD"); env != "" {
		c.DatabaseConfig.Password = env
	}
	if env := os.Getenv("DB_NAME"); env != "" {
		c.DatabaseConfig.DB = env
	}
	if env := os.Getenv("DB_PORT"); env != "" {
		if port, err := strconv.Atoi(env); err == nil {
			c.DatabaseConfig.Port = port
		} else {
			slog.Error("Invalid DB_PORT", "value", env, "error", err.Error())
		}
	}

}
