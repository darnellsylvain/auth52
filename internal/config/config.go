package config

import (
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
    Port       string        `env:"PORT" envDefault:"8080"`
    DBURL      string        `env:"AUTH52_DB_URL,required"`
    JWTSecret  string        `env:"AUTH52_JWT_SECRET,required"`
    JWTExpiry  time.Duration `env:"JWT_EXPIRY" envDefault:"1h"`
}

func Load() (*Config, error) {
    cfg := &Config{}
    if err := env.Parse(cfg); err != nil {
        return nil, err
    }
    return cfg, nil
}