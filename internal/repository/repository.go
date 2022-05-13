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
	FindByEmail(email string) *model.User
	Insert(email, hashedPassword string) *model.User
}

type TokenRepository interface {
	Find(userId uint32) *model.Token
	Insert(user *model.User, tokenVal string) *model.Token
}
