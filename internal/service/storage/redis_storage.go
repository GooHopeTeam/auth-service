package storage

import (
	"github.com/go-redis/redis"
	"time"
)

type RedisStorage struct {
	Storage
	redisClient *redis.Client
}

func NewRedisStorage(redisClient *redis.Client) Storage {
	return RedisStorage{redisClient: redisClient}
}

func (s RedisStorage) Insert(key string, data map[string]string, expiration time.Duration) error {
	genericData := make(map[string]interface{})
	for key, value := range data {
		genericData[key] = value
	}

	err := s.redisClient.HMSet(key, genericData).Err()
	if err != nil {
		return err
	}

	err = s.redisClient.Expire(key, expiration).Err()
	return err
}

func (s RedisStorage) Get(key string) (map[string]string, error) {
	data, err := s.redisClient.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s RedisStorage) Delete(key string) error {
	err := s.redisClient.Del(key).Err()
	return err
}
