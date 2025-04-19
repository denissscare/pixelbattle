package config

import (
	"flag"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host        string        `yaml:"host" env-default:"localhost"`
		Port        int           `yaml:"port" env-default:"9090"`
		Timeout     time.Duration `yaml:"timeout" env-default:"10s"`
		IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"120s"`
	} `yaml:"server"`
}

func LoadConfig() *Config {

	path := fetchConfigPath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	var config Config

	if err := viper.ReadInConfig(); err != nil {
		panic("failed to read config file: " + err.Error())
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic("failed to parse config")
	}

	return &config
}

func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}
