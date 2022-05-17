package storage

import "time"

type Storage interface {
	Insert(key string, data map[string]string, expiration time.Duration) error
	Get(key string) (map[string]string, error)
	Delete(key string) error
}
