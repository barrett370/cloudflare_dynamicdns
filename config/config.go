package config

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Cloudflare []CloudflareConfig `json:"cloudflare,omitempty"`
}

type CloudflareConfig struct {
	Authentication  CloudflareAuthentication `json:"authentication"`
	ZoneName        string                   `json:"zone_name"`
	IntervalSeconds int64                    `json:"interval_seconds"`
}

type CloudflareAuthentication struct {
	APIToken string `json:"api_token"`
}

var Config Configuration

func Load(path string) error {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(fileBytes, &Config)
	if err != nil {
		return err
	}
	return nil
}
