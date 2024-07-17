package main

import (
	"fmt"
	"io"
	"os"

	"github.com/pelletier/go-toml"
)

// Config модель конфига.
type Config struct {
	Logger  LoggerConf
	Server  ServerConf
	Metrics MetricsConf
}

// LoggerConf модель конфига логгера.
type LoggerConf struct {
	Level string `toml:"level"`
}

// ServerConf модель конфига сервера.
type ServerConf struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

// MetricsConf модель конфига метрик.
type MetricsConf struct {
	LoadAverage bool `toml:"loadAverage"`
	CPU         bool `toml:"cpu"`
	Disc        bool `toml:"disc"`
}

// NewConfig инициализация конфига.
func NewConfig(path string) (Config, error) {
	var config Config
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("open config file: %w", err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return config, fmt.Errorf("reading config file: %w", err)
	}

	err = toml.Unmarshal(b, &config)
	if err != nil {
		return config, fmt.Errorf("reading config file: %w", err)
	}

	return config, nil
}
