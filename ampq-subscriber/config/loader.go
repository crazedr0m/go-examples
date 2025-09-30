package config

import (
	"os"

	yaml "gopkg.in/yaml.v3"
)

// LoadConfig загружает конфигурацию из файла и переменных окружения
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// Читаем переменную окружения
	rabbitUrl := os.Getenv("RABBIT_URL")
	if rabbitUrl != "" {
		config.Ampq.Url = rabbitUrl
	}

	return &config, nil
}