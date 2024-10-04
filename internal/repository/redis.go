package repository

import (
	"context"
	"fmt"

	"github.com/nsvirk/moneybotstds/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	rdb *redis.Client
}

func NewRedisClient(cfg *config.Config) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{rdb: rdb}, nil
}

func (c *RedisClient) PublishTicks(channel string, tickJSON []byte) error {
	ctx := context.Background()

	err := c.rdb.Publish(ctx, channel, tickJSON).Err()
	if err != nil {
		return fmt.Errorf("failed to publish tick: %w", err)
	}

	return nil
}

func (c *RedisClient) Close() error {
	return c.rdb.Close()
}
