package repository

import (
	"fmt"
	"github.com/goohopeteam/auth-service/internal/model"
)

type Error string

func (e Error) Error() string {
	return fmt.Sprintf("repository: %s", string(e))
}

const UniqueViolationErr = Error("unique_violation")
const ModelNotFoundErr = Error("model_not_found")

type UserRepository interface {
	FindByEmail(email string) (*model.User, error)
	Insert(email, hashedPassword string) (*model.User, error)
}

type TokenRepository interface {
	Find(userId uint32) (*model.Token, error)
	Insert(user *model.User, tokenVal string) (*model.Token, error)
}
