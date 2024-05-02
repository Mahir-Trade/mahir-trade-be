package infra

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

func LoadPgDatabaseCfg() (*DatabaseCfg, error) {
	var cfg DatabaseCfg
	prefix := "PG"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}

	return &cfg, nil
}

func LoadAppCfg() (*AppCfg, error) {
	var cfg AppCfg
	prefix := "APP"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}
	return &cfg, nil
}

func LoadJwtCfg() (*JwtCfg, error) {
	var cfg JwtCfg
	prefix := "JWT"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}
	return &cfg, nil
}
