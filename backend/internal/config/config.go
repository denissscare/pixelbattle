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

	NATS struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"nats"`
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	var cfgDir string
	flag.StringVar(&cfgDir, "config", "", "path to config directory")
	flag.Parse()
	if cfgDir == "" {
		cfgDir = filepath.Join(getProjectRoot(), "internal", "config")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(cfgDir)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.BindEnv("environment", "ENVIRONMENT")
	viper.BindEnv("server.host", "SERVER_HOST")
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("redis.user", "REDIS_USER")
	viper.BindEnv("nats.url", "NATS_URL")

	if err := viper.ReadInConfig(); err != nil {
		panic("no config.yaml found")
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
