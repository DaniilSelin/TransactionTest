package config

import (
	"os"
	"log"
	"fmt"
	"time"

    	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var ConfigFilePath = "config/config.local.yml"

type ServerConfig struct {
	Host string `mapstructure:"host"` 
	Port int    `mapstructure:"port"` 
}

type MigrationConfig struct {
	Driver string `mapstructure:"driver"`
	Dir    string `mapstructure:"directory"`
}

type WalletsSeedConfig struct {
	Enabled    bool   `yaml:"Enabled"`
	FailOnError bool  `yaml:"FailOnError"`
	Count      int    `yaml:"Count"`
	Balance    float64    `yaml:"Balance"`
	MarkerFile string `yaml:"Marker_file"`	
}

type SeedingConfig struct {
    Wallets WalletsSeedConfig `yaml:"wallets"`
}

type ConnConfig struct {
	Host            string `mapstructure:"Host"`
	Port            int    `mapstructure:"Port"`
	Database        string `mapstructure:"Database"`
	User            string `mapstructure:"User"`
	Password        string `mapstructure:"Password"`
	SSLMode         string `mapstructure:"SSLMode"`
	ConnectTimeout  int `mapstructure:"ConnectTimeout"`
}

func (c *ConnConfig) ConnString() string {
    return fmt.Sprintf(
        "postgres://%s:%s@%s:%d/%s?sslmode=%s&connect_timeout=%d",
        c.User,
        c.Password,
        c.Host,
        c.Port,
        c.Database,
        c.SSLMode,
        c.ConnectTimeout,
    )
}

type PostgresPoolConfig struct {
	ConnConfig            ConnConfig `mapstructure:"ConnConfig"`
	MaxConnLifetime       time.Duration `mapstructure:"MaxConnLifetime"`
	MaxConnLifetimeJitter time.Duration `mapstructure:"MaxConnLifetimeJitter"`
	MaxConnIdleTime time.Duration `mapstructure:"MaxConnIdleTime"`
	MaxConns int32 `mapstructure:"MaxConns"`
	MinConns int32 `mapstructure:"MinConns"`
	HealthCheckPeriod time.Duration `mapstructure:"HealthCheckPeriod"`
}

type PostgresConfig struct {
	Pool PostgresPoolConfig `mapstructure:"pool"`
	ConnectRetries int `mapstructure:"ConnectRetries"` 
	ConnectRetryDelay time.Duration `mapstructure:"ConnectRetryDelay"` 
	Schema string `mapstructure:"Schema"`
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
	Migrations MigrationConfig `yaml:"migrations"`
	Seeding    SeedingConfig
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