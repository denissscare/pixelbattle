package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Environment string `mapstructure:"environment"`

	Server struct {
		Host        string        `mapstructure:"host"`
		Port        int           `mapstructure:"port"`
		Timeout     time.Duration `mapstructure:"timeout"`
		IdleTimeout time.Duration `mapstructure:"idle_timeout"`
	} `mapstructure:"server"`

	Redis struct {
		Host        string        `mapstructure:"host"`
		Port        int           `mapstructure:"port"`
		Password    string        `mapstructure:"password"`
		User        string        `mapstructure:"user"`
		DB          int           `mapstructure:"db"`
		MaxRetries  int           `mapstructure:"max_retries"`
		DialTimeout time.Duration `mapstructure:"dial_timeout"`
		Timeout     time.Duration `mapstructure:"timeout"`
	} `mapstructure:"redis"`
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found, relying on real environment variables")
	}
	envMap, _ := godotenv.Read()

	var cfgDir string
	flag.StringVar(&cfgDir, "config", "", "path to config directory")
	flag.Parse()

	if cfgDir == "" {
		cfgDir = filepath.Join(getProjectRoot(), "internal", "config")
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(cfgDir)

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic("no config.yaml found")
	}

	for rawKey, rawVal := range envMap {
		key := strings.ToLower(strings.ReplaceAll(rawKey, "_", "."))
		viper.Set(key, rawVal)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Sprintf("unable to decode into struct: %v", err))
	}

	return &cfg
}

func getProjectRoot() string {
	pwd, _ := os.Getwd()
	return pwd
}
