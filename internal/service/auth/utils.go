package auth

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/goohopeteam/auth-service/internal/model"
	"golang.org/x/crypto/sha3"
)

func makeHash(value, salt string) string {
	hasher := sha3.New512()
	hasher.Write([]byte(value + salt))
	return hex.EncodeToString(hasher.Sum(nil))
}

func generateToken(user *model.User) (string, error) {
	token := make([]byte, 64)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

func findInSlice[Key any](keys []Key, predicate func(key *Key) bool) *Key {
	for _, key := range keys {
		if predicate(&key) {
			return &key
		}
	}
	return nil
}
