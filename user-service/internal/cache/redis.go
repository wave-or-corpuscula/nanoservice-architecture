package cache

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	"userservice/pkg/utils"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	rdb *redis.Client
}

var (
	ErrCacheMiss           = errors.New("cache miss")
	ErrRedisNotInitialized = errors.New("redis is not initialized")
)

func InitRedis() (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT"),
		),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       utils.GetEnvAsInt("REDIS_DB", 0),
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to Redis:", pong)

	return &RedisClient{
		rdb: rdb,
	}, nil
}

func (rc *RedisClient) Close() error {
	if rc == nil {
		return ErrRedisNotInitialized
	}

	return rc.rdb.Close()
}

func (rc *RedisClient) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	if rc == nil {
		return ErrRedisNotInitialized
	}

	return rc.rdb.Set(ctx, key, value, expiration).Err()
}

func (rc *RedisClient) Get(ctx context.Context, key string) (string, error) {
	if rc == nil {
		return "", ErrRedisNotInitialized
	}

	value, err := rc.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return value, ErrCacheMiss
		}
		return value, err
	}

	return value, nil
}

func (rc *RedisClient) Del(ctx context.Context, keys ...string) error {
	if rc == nil {
		return ErrRedisNotInitialized
	}

	return rc.rdb.Del(ctx, keys...).Err()
}
