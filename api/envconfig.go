package api

import (
	"errors"

	"github.com/spf13/viper"
)

// EnvConfig holds env variables as fields.
type EnvConfig struct {
	GenaiAPIKey string `mapstructure:"GEMINI_API_KEY"`
	MONGODB_URI string `mapstructure:"MONGODB_DEV_URI"`
}

var (
	ErrEnvConfigFileNotFound   = errors.New("env config file not found")
	ErrFailedToReadConfig      = errors.New("failed to read in env variables")
	ErrFailedToUnmarshalConfig = errors.New("failed to unmarshal config to env config")
)

// LoadEnvConfig
func LoadEnvConfig() (EnvConfig, error) {
	viper.SetConfigName("app")
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	envConfig := EnvConfig{}
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

	return envConfig, nil
}
