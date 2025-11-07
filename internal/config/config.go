package config

import (
	"github.com/mhirii/rest-template/internal/logging"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

// Config is the global configuration struct
// Set default values ONLY here
var defaultConfig = Config{
	Port:     8888,
	LogLevel: "info",
}

type Config struct {
	Port     int    `mapstructure:"port" yaml:"port" env:"SERVICE_PORT"`
	LogLevel string `mapstructure:"log_level" yaml:"log_level" env:"LOG_LEVEL"`
	// Add more fields below, using struct tags for mapstructure/yaml/env/flag
	// Example:
	// DBHost string `mapstructure:"db_host" yaml:"db_host" env:"DB_HOST" flag:"db_host"`
	// DBPort int    `mapstructure:"db_port" yaml:"db_port" env:"DB_PORT" flag:"db_port"`
}

var (
	config Config
	loaded bool
)

// Load initializes the config ONCE, merging CLI > YAML > ENV > defaults
func Load() {
	l := logging.L()
	l.Debug().Msg("loading config")
	if loaded {
		return
	}
	loaded = true

	// Set defaults in Viper
	viper.SetDefault("port", defaultConfig.Port)
	viper.SetDefault("log_level", defaultConfig.LogLevel)

	// Bind CLI flags using pflag
	pflag.Int("port", defaultConfig.Port, "Port to listen on")
	pflag.String("log_level", defaultConfig.LogLevel, "Log level")
	pflag.String("config", "", "Path to config file")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// Load ENV variables
	viper.AutomaticEnv()
	viper.BindEnv("port", "SERVICE_PORT")
	viper.BindEnv("log_level", "LOG_LEVEL")

	// Determine config file path from CLI or ENV, default to current directory
	configPath := viper.GetString("config")
	if configPath == "" {
		configPath = viper.GetString("CONFIG_PATH")
	}
	if configPath == "" {
		configPath = "."
	}
	// If configPath is a file, use SetConfigFile; if directory, use AddConfigPath
	if fi, err := os.Stat(configPath); err == nil && !fi.IsDir() {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(configPath)
	}
	_ = viper.ReadInConfig() // ignore error if file not found

	// Unmarshal to config struct
	config = defaultConfig // start with defaults
	_ = viper.Unmarshal(&config)
}

// Get returns the global config
func Get() Config {
	return config
}
