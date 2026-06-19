package config

import (
	"encoding/json"
	"os"

	"github.com/SayukiDev/VRCLotterySystem/internal/global/validator"
)

type Config struct {
	LogLevel        string     `json:"log_level" validate:"oneof=debug info warn error fatal"`
	DataPath        string     `json:"data_path" validate:"required"`
	DiscordMasterID string     `json:"discord_master_id" validate:"required"`
	Terms           string     `json:"terms" validate:"required"`
	Form            []FormItem `json:"form" validate:"required"`
}

type FormItem struct {
	IsId     bool     `json:"is_id" validate:"required"`
	Title    string   `json:"title" validate:"required"`
	Desc     string   `json:"desc"`
	Required bool     `json:"required"`
	Options  []string `json:"options"`
	Type     string   `json:"type" validate:"required oneof=content input options"`
}

func DefaultConfig() *Config {
	return &Config{
		LogLevel:        "info",
		DataPath:        "./data",
		DiscordMasterID: "Input your discord id here",
		Form: []FormItem{
			{
				Title: "Input your item title here",
				Desc:  "Input your item description here",
				Type:  "content",
			},
			{
				Title:    "Input your item title here",
				Required: true,
				Type:     "input",
			},
			{
				Title:    "Input your item title here",
				Required: true,
				Type:     "options",
			},
		},
	}
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	c := DefaultConfig()
	err = json.NewDecoder(f).Decode(c)
	if err != nil {
		return nil, err
	}
	err = c.Validate()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) Validate() error {
	return validator.V.Struct(c)
}
