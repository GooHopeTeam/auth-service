package repository

import (
	"database/sql"
	"time"

	"github.com/goohopeteam/auth-service/internal/model"
	"github.com/jmoiron/sqlx"
)

type UserRepositoryImpl struct {
	UserRepository
	db *sqlx.DB
}

type TokenRepositoryImpl struct {
	TokenRepository
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return UserRepositoryImpl{db: db}
}

func NewTokenRepository(db *sqlx.DB) TokenRepository {
	return TokenRepositoryImpl{db: db}
}

func (rep UserRepositoryImpl) FindByEmail(email string) *model.User {
	var user model.User
	err := rep.db.Get(&user, "SELECT id, email, hashed_password, created_at FROM \"user\" WHERE email = $1", email)
	if err != nil && err == sql.ErrNoRows {
		return nil
	}
	return &user
}

func (rep UserRepositoryImpl) Insert(email, hashedPassword string) *model.User {
	var user model.User
	err := rep.db.Get(&user, "INSERT INTO \"user\"(email, hashed_password, created_at) VALUES($1, $2, $3) RETURNING *;",
		email, hashedPassword, time.Now())
	if err != nil {
		return nil
	}
	return &user
}

func (rep TokenRepositoryImpl) Find(userId uint32) *model.Token {
	var token model.Token
	err := rep.db.Get(&token, "SELECT token FROM token WHERE user_id = $1", userId)
	if err != nil {
		return nil
	}
	return &token
}

func (rep TokenRepositoryImpl) Insert(user *model.User, tokenVal string) *model.Token {
	var token model.Token
	err := rep.db.Get(&token, "INSERT INTO \"token\"(user_id, token) VALUES($1, $2) RETURNING *;",
		user.Id, tokenVal)
	if err != nil {
		return nil
	}
	return &token
}
