package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Admins       []string `yaml:"admins"`
	WhisperToken string   `yaml:"token"`
	WhisperURI   string   `yaml:"whisperURI"`
	DeviceID     string   `yaml:"deviceID"`
}

func ReadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
