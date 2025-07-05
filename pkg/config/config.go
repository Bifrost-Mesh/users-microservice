package config

type (
	Config struct {
		DevMode      bool `yaml:"devMode"      default:"False"`
		DebugLogging bool `yaml:"debugLogging" default:"False"`

		ServerPort int `yaml:"serverPort" default:"4000" validate:"gt=0"`

		JWTSigningKey string `yaml:"jwtSigningKey" validate:"notblank"`

		Postgres PostgresConfig `yaml:"postgres" validate:"required"`
	}

	PostgresConfig struct {
		URL string `yaml:"url" validate:"notblank"`
	}
)
