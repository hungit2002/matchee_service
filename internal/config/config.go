package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	AppName         string
	HTTPHost        string
	HTTPPort        string
	MySQLDSN        string
	RedisAddr       string
	RedisDB         int
	RedisPassword   string
	RabbitURL       string
	JWTSecret       string
	JWTExpiry       time.Duration
	RefreshExpiry   time.Duration
	ShutdownTimeout time.Duration
}

func Load() (*Config, error) {
	c := &Config{
		AppName:         getEnv("APP_NAME", "matchee-services"),
		HTTPHost:        getEnv("HTTP_HOST", "0.0.0.0"),
		HTTPPort:        getEnv("HTTP_PORT", "8080"),
		MySQLDSN:        getEnv("MYSQL_DSN", "root:root@tcp(127.0.0.1:3306)/matchee?charset=utf8mb4&parseTime=True&loc=Local"),
		RedisAddr:       getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisDB:         getEnvInt("REDIS_DB", 0),
		RedisPassword:   os.Getenv("REDIS_PASSWORD"),
		RabbitURL:       getEnv("RABBITMQ_URL", "amqp://guest:guest@127.0.0.1:5672/"),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpiry:       getEnvDuration("JWT_EXPIRY", 15*time.Minute),
		RefreshExpiry:   getEnvDuration("REFRESH_EXPIRY", 7*24*time.Hour),
		ShutdownTimeout: getEnvDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
	}
	return c, nil
}

func (c *Config) HTTPAddr() string {
	return fmt.Sprintf("%s:%s", c.HTTPHost, c.HTTPPort)
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		var out int
		_, _ = fmt.Sscanf(v, "%d", &out)
		if out != 0 || v == "0" {
			return out
		}
	}
	return def
}

func getEnvDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
