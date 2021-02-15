package main

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// config describes the available configuration
// of the running service
type config struct {
	Debug             bool
	MaxExpirationDays *int `mapstructure:"max_expiration_days"`
}

// Validate makes sure that the config makes sense
func (c *config) Validate() error {
	if &c.MaxExpirationDays == nil {
		return errors.New("Config: expiration time is required")
	}
	return nil
}

// Set the file name of the configurations file
func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("janitor")

	defaults := map[string]interface{}{
		"debug":               false,
		"environment":         "dev",
		"max_expiration_days": 432,
	}
	for key, value := range defaults {
		viper.SetDefault(key, value)
	}
}

// LoadConfig checks file and environment variables
func LoadConfig(logger log.FieldLogger) error {
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return errors.Wrap(err, "config load")
	}
	return errors.Wrap(cfg.Validate(), "config validate")
}
