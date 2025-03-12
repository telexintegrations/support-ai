package api

import (
	"errors"

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
	viper.SetConfigType("env")
	viper.AutomaticEnv() // Always read from environment variables

	// Check if running locally by looking for a .env file
	viper.SetConfigFile(".env") // Explicitly set .env file
	if err := viper.ReadInConfig(); err != nil {
		// If the file isn't found, just continue using environment variables
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return EnvConfig{}, ErrFailedToReadConfig
		}
	}

	var envConfig EnvConfig
	if err := viper.Unmarshal(&envConfig); err != nil {
		return EnvConfig{}, ErrFailedToUnmarshalConfig
	}

	return envConfig, nil
}
