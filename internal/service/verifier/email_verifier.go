package verifier

import (
	"fmt"
	"github.com/goohopeteam/auth-service/internal/service/mailer"
	"github.com/goohopeteam/auth-service/internal/service/storage"
	"time"
)

type EmailVerifier struct {
	mailer       mailer.Mailer
	storage      storage.Storage
	codeProvider codeProvider
}

func NewEmailVerifier(mailer mailer.Mailer, storage storage.Storage) Verifier {
	return EmailVerifier{mailer: mailer, storage: storage, codeProvider: newCodeProvider()}
}

func (v EmailVerifier) Send(email string, data map[string]string) error {
	code := v.codeProvider.generateCode()
	message := fmt.Sprintf("Hello, your verification code is %s", code)

	err := v.mailer.Send(email, message)
	if err != nil {
		return err
	}

	data["__code"] = code
	err = v.storage.Insert(email, data, 10*time.Minute)
	return err
}

func (v EmailVerifier) Check(email, code string) (map[string]string, error) {
	data, err := v.storage.Get(email)
	if data == nil || err != nil {
		return nil, err
	}

	if data["__code"] != code {
		return nil, nil
	}

	err = v.storage.Delete(email)
	if err != nil {
		return nil, err
	}
	delete(data, "__code")

	return data, nil
}
