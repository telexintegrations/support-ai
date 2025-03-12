package api

import (
	"errors"
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
	envConfig := EnvConfig{}
	if os.Getenv("NODE_ENV") == "production" {
        // Production environment logic
		viper.AutomaticEnv()
    } else {
        // Development or other environment logic
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return envConfig, ErrEnvConfigFileNotFound
		}
		return envConfig, ErrFailedToReadConfig
	}

	err = viper.Unmarshal(&envConfig)
	if err != nil {
		return envConfig, ErrFailedToUnmarshalConfig
	}
    }

	return envConfig, nil
}
