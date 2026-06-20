package config

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/SayukiDev/VRCLotterySystem/internal/global/validator"
)

type Config struct {
	LogLevel        string     `json:"log_level" validate:"oneof=debug info warn error fatal"`
	DataPath        string     `json:"data_path" validate:"required"`
	DiscordMasterID string     `json:"discord_master_id" validate:"required"`
	SiteData        SiteData   `json:"site_data" validate:"required"`
	Form            []FormItem `json:"form" validate:"required"`
}

type SiteData struct {
	Title     string       `json:"title" validate:"required"`      // A title for the site
	FormTitle string       `json:"form_title" validate:"required"` // A title for the form page
	Terms     StringOrFile `json:"terms" validate:"required"`      // Markdown
}

type FormItem struct {
	IsId     bool         `json:"is_id"`
	Title    string       `json:"title" validate:"required"`
	Desc     StringOrFile `json:"desc"`
	Required bool         `json:"required"`
	Options  []string     `json:"options"`
	Type     string       `json:"type" validate:"required,oneof=content input options"`
}

type StringOrFile string

func (s *StringOrFile) UnmarshalJSON(b []byte) error {
	body := ""
	err := json.Unmarshal(b, &body)
	if err != nil {
		return err
	}
	if body == "" {
		return nil
	}
	if strings.HasSuffix(body, ".md") {
		bs, err := os.ReadFile(body)
		if err != nil {
			return err
		}
		*s = StringOrFile(bs)
	} else {
		*s = StringOrFile(body)
	}
	return nil
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
	defer f.Close()
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
	haveID := false
	for _, v := range c.Form {
		if v.IsId {
			haveID = true
		}
	}
	if !haveID {
		return errors.New("no id found in form")
	}
	return validator.V.Struct(c)
}
