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

	Postgres struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"postgres"`

	Minio struct {
		Endpoint   string `mapstructure:"endpoint"`
		AccessKey  string `mapstructure:"access_key"`
		SecretKey  string `mapstructure:"secret_key"`
		Bucket     string `mapstructure:"bucket"`
		UseSSL     bool   `mapstructure:"use_ssl"`
		PublicHost string `mapstructure:"public"`
	} `mapstructure:"minio"`
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

	viper.BindEnv("postgres.host", "POSTGRES_HOST")
	viper.BindEnv("postgres.port", "POSTGRES_PORT")
	viper.BindEnv("postgres.user", "POSTGRES_USER")
	viper.BindEnv("postgres.password", "POSTGRES_PASSWORD")
	viper.BindEnv("postgres.dbname", "POSTGRES_DB")
	viper.BindEnv("postgres.sslmode", "POSTGRES_SSLMODE")

	viper.BindEnv("minio.endpoint", "MINIO_ENDPOINT")
	viper.BindEnv("minio.public", "MINIO_PUBLIC_HOST")
	viper.BindEnv("minio.access_key", "MINIO_ACCESS_KEY")
	viper.BindEnv("minio.secret_key", "MINIO_SECRET_KEY")
	viper.BindEnv("minio.bucket", "MINIO_BUCKET")
	viper.BindEnv("minio.use_ssl", "MINIO_USE_SSL")

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
