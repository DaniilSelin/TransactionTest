package config

import (
	"os"
	"log"
	"fmt"

	"gopkg.in/yaml.v2"
	"go.uber.org/zap"
)

var ConfigFilePath = "config/config.local.yml"

type ServerConfig struct {
	Host string `yaml:"SERVER_HOST" env:"SERVER_HOST" env-default:"0.0.0.0"` 
	Port int    `yaml:"SERVER_PORT" env:"SERVER_PORT" env-default:"8080"` 
}

type DatabaseConfig struct {
	Host          string `yaml:"POSTGRES_HOST" env:"POSTGRES_HOST" env-default:"localhost"`
	Port          int    `yaml:"POSTGRES_PORT" env:"POSTGRES_PORT" env-default:"5432"`
	User          string `yaml:"POSTGRES_USER" env:"POSTGRES_USER" env-default:"postgres"`
	Password      string `yaml:"POSTGRES_PASS" env:"POSTGRES_PASS" env-default:"changeme"`
	Dbname        string `yaml:"POSTGRES_DB" env:"POSTGRES_DB" env-default:"transaction_system"`
	Sslmode       string `yaml:"POSTGRES_SslMODE" env:"POSTGRES_SslMODE" env-default:"disable"`
	Schema		  string `yaml:"POSTGRES_SCHEMA" env:"POSTGRES_SCHEMA" env-default:"TransactionSystem"`
}

type LoggerConfig struct {
	Logger zap.Config `yaml:",inline"`
}

func (l *LoggerConfig) Build() (*zap.Logger, error) {
	return l.Logger.Build()
}

type Config struct {
	Database DatabaseConfig `yaml:"postgres"`
	Server   ServerConfig   `yaml:"server"`
	LoggerConfig   LoggerConfig   `yaml:"logger"`
}

func LoadConfig() (Config, error) {
	cfgFile := os.Getenv("CONFIG_FILE_PATH"); 
	if cfgFile == "" {
		cfgFile = ConfigFilePath
		log.Println("WARN:failed to read CONFIG_FILE_PATH, using default path")
	}

	file, err := os.Open(cfgFile)
	if err != nil {
		return Config{}, fmt.Errorf("FATAL: error reading config file: %w", err)
	}
	defer file.Close()

	var cfg Config

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("FATAL:unable to decode into struct: %w", err)
	}

	return cfg, nil
}