package verifier

import (
	"github.com/goohopeteam/auth-service/internal/service/mailer"
	"github.com/goohopeteam/auth-service/internal/service/storage"
	"time"
)

type MailerMock struct {
	mailer.Mailer
}

func newMailerMock() mailer.Mailer {
	return &MailerMock{}
}

func (m MailerMock) Send(email, message string) error {
	return nil
}

type StorageMock struct {
	storage.Storage
	actualStorage map[string]map[string]string
}

func newStorageMock() *StorageMock {
	return &StorageMock{actualStorage: make(map[string]map[string]string)}
}

func (s *StorageMock) Insert(key string, data map[string]string, expiration time.Duration) error {
	s.actualStorage[key] = data
	return nil
}

func (s *StorageMock) Get(key string) (map[string]string, error) {
	return s.actualStorage[key], nil
}

func (s *StorageMock) Delete(key string) error {
	delete(s.actualStorage, key)
	return nil
}

type codeProviderMock struct {
	codeProvider
}

func newCodeProviderMock() codeProviderMock {
	return codeProviderMock{}
}

func (p codeProviderMock) generateCode() string {
	return "1"
}
