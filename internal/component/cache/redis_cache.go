package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"fiangumilar.id/e-wallet/domain"
	"fiangumilar.id/e-wallet/internal/config"
	"github.com/redis/go-redis/v9"
)

type redisCacheRepository struct {
	rdb *redis.Client
}

func NewRedisClient(conf *config.Config) (domain.CacheRepository, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Pass,
		DB:       0,
	})

	// Test the connection with Redis using PING
	ctx := context.Background()
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("Ping Error: %v", err)
	}

	fmt.Println("Connected to Redis successfully:", pong)

	return &redisCacheRepository{
		rdb: rdb,
	}, nil
}

// Get implements domain.CacheRepository.
func (r redisCacheRepository) Get(key string) ([]byte, error) {
	val, err := r.rdb.Get(context.Background(), key).Result()
	if err != nil {
		log.Printf("Failed get token %s", err)
		return nil, err
	}
	log.Println("Successfully retrieved tokens")
	return []byte(val), nil
}

// Set implements domain.CacheRepository.
func (r redisCacheRepository) Set(key string, entry []byte) error {
	return r.rdb.Set(context.Background(), key, entry, 15*time.Minute).Err()
}
