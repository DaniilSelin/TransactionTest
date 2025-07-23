package config

import (
	"os"
	"log"
	"fmt"
	"time"
    "reflect"

    "github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ConfigFilePath = "config/config.local.yml"

type ServerConfig struct {
	Host string `mapstructure:"host"` 
	Port int    `mapstructure:"port"` 
}

type MigrationConfig struct {
	ConnectRetries int `mapstructure:"ConnectRetries"` 
	ConnectRetryDelay time.Duration `mapstructure:"ConnectRetryDelay"` 
}

type PostgresConfig struct {
	Pool pgxpool.Config `mapstructure:"postgres"`
}

type LoggerConfig struct {
	Logger zap.Config `mapstructure:"logger"`
}

func (l *LoggerConfig) Build() (*zap.Logger, error) {
	return l.Logger.Build()
}

type Config struct {
	Postgres PostgresConfig `yaml:"postgres"`
	Server   ServerConfig   `yaml:"server"`
	Logger   LoggerConfig   `yaml:"logger"`
	Migration MigrationConfig `yaml:"migration"`
}

// zapLevelHook: строка в zap.AtomicLevel
func zapLevelHook(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
    if to != reflect.TypeOf(zap.AtomicLevel{}) {
        return data, nil
    }

    s, ok := data.(string)
    if !ok {
        return data, nil
    }

    var lvl zap.AtomicLevel
    if err := lvl.UnmarshalText([]byte(s)); err != nil {
        return nil, err
    }
    return lvl, nil
}

func LoadConfig() (Config, error) {
	if cfgFile := os.Getenv("CONFIG_FILE_PATH"); cfgFile != "" {
	    viper.SetConfigFile(cfgFile)
	} else {
	    viper.SetConfigFile(ConfigFilePath)
		log.Println("WARN:failed to read CONFIG_FILE_PATH, using default path")
	}

	if err := viper.ReadInConfig(); err != nil {
	    return Config{}, fmt.Errorf("FATAL: error reading config file: %w", err)
	}

	var cfg Config

	decoderConfig := &mapstructure.DecoderConfig{
        DecodeHook: mapstructure.ComposeDecodeHookFunc(
            mapstructure.StringToTimeDurationHookFunc(), // чтобы мапить durations
            zapLevelHook,                                // хук для AtomicLevel
        ),
        Result:  &cfg,
        TagName: "mapstructure",
    }

    dec, err := mapstructure.NewDecoder(decoderConfig)
    if err != nil {
        return Config{}, fmt.Errorf("FATAL:unable to create new decoder: %w", err)
    }

    if err := dec.Decode(viper.AllSettings()); err != nil {
        return Config{}, fmt.Errorf("FATAL:unable to decode into struct: %w", err)
    }

	return cfg, nil
}