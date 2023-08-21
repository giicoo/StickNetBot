package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	TOKEN string `env:"BOT_TOKEN"`
}

func InitConfig(path string) (Config, error) {
	cfg := Config{}
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("ConfigErr: %v", err)
	}
	return cfg, nil
}
