package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type LimiterConfig struct {
	Enabled bool    `env:"LIMITER_ENABLED" envDefault:"true"`
	RPS     float64 `env:"LIMITER_RPS" envDefault:"5"`
	Burst   int     `env:"LIMITER_BURST" envDefault:"10"`
}

type Config struct {
	Port      string        `env:"PORT" envDefault:"8080"`
	DBURL     string        `env:"AUTH52_DB_URL,required"`
	JWTSecret string        `env:"AUTH52_JWT_SECRET,required"`
	JWTExpiry time.Duration `env:"JWT_EXPIRY" envDefault:"1h"`
	Limiter   LimiterConfig
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
