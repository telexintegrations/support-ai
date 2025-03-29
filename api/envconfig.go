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
	CHROMADB_DEV_URI   string `mapstructure:"CHROMADB_DEV_URI"`
	TELEX_URI          string `mapstructure:"TELEX_URI"`
}

var (
	ErrEnvConfigFileNotFound   = errors.New("env config file not found")
	ErrFailedToReadConfig      = errors.New("failed to read in env variables")
	ErrFailedToUnmarshalConfig = errors.New("failed to unmarshal config to env config")
)

// LoadEnvConfig
func LoadEnvConfig() (EnvConfig, error) {
	var envConfig EnvConfig

	// Always load system environment variables
	viper.AutomaticEnv()

	// Check if we're running in production
	isProduction := os.Getenv("NODE_ENV") == "production"

	if !isProduction {
		// Load .env file in non-production environments
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		viper.AddConfigPath(".")

		// Attempt to read config file
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				fmt.Println("❌ Failed to read config file:", err)
				return envConfig, fmt.Errorf("failed to read config: %w", err)
			}
			fmt.Println("⚠️  No .env file found, falling back to system environment variables")
		}
	}

	// Unmarshal environment variables into struct
	if err := viper.Unmarshal(&envConfig); err != nil {
		fmt.Println("❌ Failed to unmarshal environment variables:", err)
		return envConfig, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if envConfig == (EnvConfig{}) {
		fmt.Println("envConfig is empty!, loading os variables")
		apikey := os.Getenv("GEMINI_API_KEY")
		uri := os.Getenv("MONGODB_DEV_URI")
		db_username := os.Getenv("MONGO_USERNAME")
		db_password := os.Getenv("MONGO_PASSWORD")
		db_name := os.Getenv("MONGODATABASE_NAME")
		chromadb_dev_uri := os.Getenv("CHROMADB_DEV_URI")
		telex_uri := os.Getenv("TELEX_URI")
		envConfig = EnvConfig{
			GenaiAPIKey:        apikey,
			MONGODB_DEV_URI:    uri,
			MONGO_USERNAME:     db_username,
			MONGO_PASSWORD:     db_password,
			MONGODATABASE_NAME: db_name,
			CHROMADB_DEV_URI:   chromadb_dev_uri,
			TELEX_URI:          telex_uri,
		}
	}

	return envConfig, nil
}
