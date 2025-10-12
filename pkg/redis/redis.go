package redis

import (
	"context"

	"matchee/services/internal/config"

	redislib "github.com/redis/go-redis/v9"
)

func NewClient(cfg *config.Config) (*redislib.Client, error) {
	client := redislib.NewClient(&redislib.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return client, nil
}
