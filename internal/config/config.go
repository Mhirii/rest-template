package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/mhirii/rest-template/internal/logging"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// --- Domain Sub-Structs ---
type ServerConfig struct {
	Port     int    `flag:"port" env:"SERVICE_PORT" yaml:"port" default:"8888" validate:"min=1,max=65535"`
	LogLevel string `flag:"log_level" env:"LOG_LEVEL" yaml:"log_level" default:"info" validate:"oneof=debug info warn error"`
}

type DBConfig struct {
	Host string `flag:"db_host" env:"DB_HOST" yaml:"db_host" default:"localhost" validate:"required"`
	Port int    `flag:"db_port" env:"DB_PORT" yaml:"db_port" default:"5432" validate:"min=1,max=65535"`
}

// --- Main Config Struct ---
type Config struct {
	Server ServerConfig
	DB     DBConfig
}

var (
	config Config
	loaded bool
)

// --- Auto-Binding and Validation Helpers ---
func bindConfigStruct(v *viper.Viper, s interface{}, prefix string) {
	t := reflect.TypeOf(s)
	vStruct := reflect.ValueOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		vStruct = vStruct.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fieldName := f.Name
		flagTag := f.Tag.Get("flag")
		envTag := f.Tag.Get("env")
		yamlTag := f.Tag.Get("yaml")
		defaultTag := f.Tag.Get("default")

		// Compose key with prefix for nested structs
		key := yamlTag
		if prefix != "" {
			key = prefix + "." + yamlTag
		}

		// Register flag
		if flagTag != "" {
			switch f.Type.Kind() {
			case reflect.Int, reflect.Int64:
				defVal := 0
				if defaultTag != "" {
					if v, err := strconv.Atoi(defaultTag); err == nil {
						defVal = v
					}
				}
				pflag.Int(flagTag, defVal, fmt.Sprintf("%s (default: %s)", fieldName, defaultTag))
			case reflect.String:
				defVal := ""
				if defaultTag != "" {
					defVal = defaultTag
				}
				pflag.String(flagTag, defVal, fmt.Sprintf("%s (default: %s)", fieldName, defaultTag))
			case reflect.Bool:
				defVal := false
				if defaultTag != "" {
					defVal = defaultTag == "true"
				}
				pflag.Bool(flagTag, defVal, fmt.Sprintf("%s (default: %s)", fieldName, defaultTag))
			case reflect.Float64:
				defVal := 0.0
				if defaultTag != "" {
					if v, err := strconv.ParseFloat(defaultTag, 64); err == nil {
						defVal = v
					}
				}
				pflag.Float64(flagTag, defVal, fmt.Sprintf("%s (default: %s)", fieldName, defaultTag))
			default:
				fmt.Fprintf(os.Stderr, "Unsupported flag type for %s: %s\n", fieldName, f.Type.Kind())
			}
		}
		// Set default
		if defaultTag != "" {
			v.SetDefault(key, defaultTag)
		}
		// Bind env
		if envTag != "" {
			v.BindEnv(key, envTag)
		}
	}
}

func validateConfigStruct(s interface{}) error {
	t := reflect.TypeOf(s)
	vStruct := reflect.ValueOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		vStruct = vStruct.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		val := vStruct.Field(i)
		// Recursively validate nested structs
		if f.Type.Kind() == reflect.Struct {
			if err := validateConfigStruct(val.Interface()); err != nil {
				return fmt.Errorf("%s: %w", f.Name, err)
			}
			continue
		}
		validateTag := f.Tag.Get("validate")
		if validateTag == "" {
			continue
		}
		for _, rule := range strings.Split(validateTag, ",") {
			rule = strings.TrimSpace(rule)
			if rule == "required" && isZero(val) {
				return fmt.Errorf("%s is required", f.Name)
			}
			if strings.HasPrefix(rule, "min=") {
				min := atoi(strings.TrimPrefix(rule, "min="))
				if val.Int() < int64(min) {
					return fmt.Errorf("%s must be >= %d", f.Name, min)
				}
			}
			if strings.HasPrefix(rule, "max=") {
				max := atoi(strings.TrimPrefix(rule, "max="))
				if val.Int() > int64(max) {
					return fmt.Errorf("%s must be <= %d", f.Name, max)
				}
			}
			if strings.HasPrefix(rule, "oneof=") {
				opts := strings.Split(strings.TrimPrefix(rule, "oneof="), " ")
				found := false
				for _, opt := range opts {
					if val.String() == opt {
						found = true
					}
				}
				if !found {
					return fmt.Errorf("%s must be one of %v", f.Name, opts)
				}
			}
		}
	}
	return nil
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int64:
		return v.Int() == 0
	}
	return false
}

func atoi(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}

// --- Loader ---
func Load() {
	l := logging.L()
	l.Debug().Msg("loading config")
	if loaded {
		return
	}
	loaded = true

	v := viper.New()

	// Bind all config fields recursively
	bindConfigStruct(v, &config.Server, "server")
	bindConfigStruct(v, &config.DB, "db")

	// Bind CLI flags
	pflag.String("config", "", "Path to config file or directory")
	pflag.Parse()
	v.BindPFlags(pflag.CommandLine)

	// Load ENV variables
	v.AutomaticEnv()

	// Determine config file path from CLI or ENV, default to current directory
	configPath := v.GetString("config")
	if configPath == "" {
		configPath = v.GetString("CONFIG_PATH")
	}
	if configPath == "" {
		configPath = "."
	}
	if fi, err := os.Stat(configPath); err == nil && !fi.IsDir() {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(configPath)
	}
	_ = v.ReadInConfig() // ignore error if file not found

	// Unmarshal to config struct
	_ = v.Unmarshal(&config)

	// Validate config
	if err := validateConfigStruct(&config.Server); err != nil {
		panic(fmt.Sprintf("Config validation error: %v", err))
	}
	if err := validateConfigStruct(&config.DB); err != nil {
		panic(fmt.Sprintf("Config validation error: %v", err))
	}
}

// Get returns the global config
func Get() Config {
	return config
}
