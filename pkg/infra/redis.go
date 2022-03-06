package infra

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisRepository struct {
	rdb *redis.Client
}

var ctx = context.Background()

func NewRedisRepository() *RedisRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &RedisRepository{rdb: rdb}
}

type CachedResponseDTO struct {
	Status int    `redis:"status"`
	Data   string `redis:"data"`
}

func (repo *RedisRepository) Get(key string) (*CachedResponseDTO, error) {
	cached := repo.rdb.HGetAll(ctx, key)
	if cached.Err() != nil {
		return nil, cached.Err()
	}

	var response CachedResponseDTO
	if err := cached.Scan(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (repo *RedisRepository) Set(key string, statusCode int, body string) error {
	_, err := repo.rdb.HSet(ctx, key, "status", statusCode, "data", body).Result()
	ttl := time.Hour
	repo.rdb.Expire(ctx, key, ttl)

	return err
}
