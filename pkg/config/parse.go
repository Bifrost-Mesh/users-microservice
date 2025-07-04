package config

import (
	"context"
	"log/slog"
	"os"

	"github.com/Bifrost-Mesh/users-microservice/pkg/assert"
	"github.com/go-playground/validator/v10"
	"github.com/mcuadros/go-defaults"
	"gopkg.in/yaml.v3"
)

// Parses and validates the config file at the given path.
// The parsed config is then returned.
//
// Panics if any error occurs.
func MustParseConfigFile(ctx context.Context,
	configFilePath string,
	validator *validator.Validate,
) *Config {
	configFileContents, err := os.ReadFile(configFilePath)
	assert.AssertErrNil(ctx, err, "Failed reading config file", slog.String("path", configFilePath))

	return MustParseConfig(ctx, configFileContents, validator)
}

// Parses and validates the given unmarshalled config.
// The parsed config is then returned.
//
// Panics if any error occurs.
func MustParseConfig(ctx context.Context,
	unparsedConfig []byte,
	validator *validator.Validate,
) *Config {
	config := new(Config)

	err := yaml.Unmarshal(unparsedConfig, config)
	assert.AssertErrNil(ctx, err, "Failed YAML unmarshalling config")

	// Validate based on struct tags.
	err = validator.Struct(config)
	assert.AssertErrNil(ctx, err, "Config validation failed")

	// Populate optional fields with corresponding default values.
	defaults.SetDefaults(config)

	return config
}
