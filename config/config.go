package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	TOKEN string `env:"BOT_TOKEN"`

	// Welcome _______________________
	WelcomeGroup struct {
		Phrase         string `yaml:"phrase"`
		InlineKeyboard struct {
			Row1 struct {
				Btn1 string `yaml:"btn_1"`
			} `yaml:"row_1"`
			Row2 struct {
				Btn1 string `yaml:"btn_1"`
			} `yaml:"row_2"`
		} `yaml:"inline_keyboard"`
	} `yaml:"welcome"`

	// ClassicFsm _____________________
	ClassicFsmGroup struct {
		Start struct {
			Phrase         string `yaml:"phrase"`
			InlineKeyboard struct {
				Row1 struct {
					Btn1 string `yaml:"btn_1"`
				} `yaml:"row_1"`
			} `yaml:"inline_keyboard"`
		} `yaml:"start"`
	} `yaml:"classicFsm"`
}

func InitConfig(path string) (Config, error) {
	cfg := Config{}
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("ConfigErr: %v", err)
	}
	return cfg, nil
}
