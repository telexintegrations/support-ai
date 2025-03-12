package api

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// EnvConfig holds env variables as fields.
type EnvConfig struct {
	GenaiAPIKey        string `mapstructure:"GEMINI_API_KEY"`
	MONGODB_DEV_URI    string `mapstructure:"MONGODB_DEV_URI"` // for development
	MONGO_USERNAME     string `mapstructure:"MONGO_USERNAME"`
	MONGO_PASSWORD     string `mapstructure:"MONGO_PASSWORD"`
	MONGODATABASE_NAME string `mapstructure:"MONGODATABASE_NAME"`
}

var (
	ErrEnvConfigFileNotFound   = errors.New("env config file not found")
	ErrFailedToReadConfig      = errors.New("failed to read in env variables")
	ErrFailedToUnmarshalConfig = errors.New("failed to unmarshal config to env config")
)

// LoadEnvConfig
func LoadEnvConfig() (EnvConfig, error) {
	var envConfig EnvConfig

	// Always use viper.AutomaticEnv() to pick up env vars
	viper.AutomaticEnv()

	// If NOT in production, load from .env file
	if os.Getenv("NODE_ENV") != "production" {
		viper.SetConfigName(".env")  // Looks for "env.env" or ".env"
		viper.SetConfigType("env")  // Explicitly define as env file
		viper.AddConfigPath(".")    // Look in the current directory

		// Read the config file if it exists
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return envConfig, fmt.Errorf("failed to read config: %w", err)
			}
			// If file not found, continue (use system env vars)
		}
	}

	// Unmarshal environment variables into struct
	if err := viper.Unmarshal(&envConfig); err != nil {
		return envConfig, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return envConfig, nil
}