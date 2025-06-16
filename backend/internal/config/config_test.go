package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	configContent := `
environment: test
server:
  host: "127.0.0.1"
  port: 8080
  timeout: 5s
  idle_timeout: 60s
redis:
  host: "localhost"
  port: 6379
  password: "test"
  user: "user"
  db: 1
  max_retries: 3
  dial_timeout: 2s
  timeout: 1s
nats:
  url: "nats://localhost:4222"
postgres:
  host: "localhost"
  port: 5432
  user: "pg"
  password: "secret"
  dbname: "pixelbattle"
  sslmode: "disable"
minio:
  endpoint: "localhost:9000"
  access_key: "minio"
  secret_key: "minio_secret"
  bucket: "test-bucket"
  use_ssl: false
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"test", "-config=" + dir}

	cfg := LoadConfig()

	require.Equal(t, "test", cfg.Environment)
	require.Equal(t, "127.0.0.1", cfg.Server.Host)
	require.Equal(t, 8080, cfg.Server.Port)
	require.Equal(t, "localhost", cfg.Redis.Host)
	require.Equal(t, 6379, cfg.Redis.Port)
	require.Equal(t, "nats://localhost:4222", cfg.NATS.URL)
	require.Equal(t, "pixelbattle", cfg.Postgres.DBName)
	require.Equal(t, "test-bucket", cfg.Minio.Bucket)
}
