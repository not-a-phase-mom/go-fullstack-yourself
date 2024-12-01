package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Redis struct {
	client *redis.Client
}

func InitRedis(addr, password string, db int) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test the connection
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	log.Println("Redis connection initialized successfully.")
	return &Redis{client: client}, nil
}

func GenerateToken(key string) string {
	// use a library to generate a random token
	// like github.com/google/uuid
	randomID := uuid.New().String()
	return randomID
}

func (r *Redis) SetValue(key string, value interface{}) error {
	if err := r.client.Set(context.Background(), key, value, 0).Err(); err != nil {
		return fmt.Errorf("failed to set value: %w", err)
	}
	return nil
}

func (r *Redis) GetValue(key string) (string, error) {
	value, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get value: %w", err)
	}
	return value, nil
}

func (r *Redis) DeleteValue(key string) error {
	if err := r.client.Del(context.Background(), key).Err(); err != nil {
		return fmt.Errorf("failed to delete value: %w", err)
	}
	return nil
}

func (r *Redis) Close() error {
	if err := r.client.Close(); err != nil {
		return fmt.Errorf("failed to close redis connection: %w", err)
	}
	return nil
}
