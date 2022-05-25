package repository

import (
	"database/sql"
	"errors"
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

func (rep UserRepositoryImpl) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := rep.db.Get(&user, "SELECT id, email, hashed_password, created_at FROM \"user\" WHERE email = $1", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (rep UserRepositoryImpl) Insert(email, hashedPassword string) (*model.User, error) {
	var user model.User
	err := rep.db.Get(&user, "INSERT INTO \"user\"(email, hashed_password, created_at) VALUES($1, $2, $3) RETURNING *;",
		email, hashedPassword, time.Now())
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (rep UserRepositoryImpl) UpdatePassword(userID uint32, hashedPassword string) (*model.User, error) {
	var user model.User
	err := rep.db.Get(&user, "UPDATE \"user\" SET hashed_password = $1 WHERE id = $2 RETURNING *;",
		hashedPassword, userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (rep TokenRepositoryImpl) Find(userID uint32) (*model.Token, error) {
	var token model.Token
	err := rep.db.Get(&token, "SELECT token FROM token WHERE user_id = $1", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (rep TokenRepositoryImpl) Insert(userID uint32, tokenVal string) (*model.Token, error) {
	var token model.Token
	err := rep.db.Get(&token, "INSERT INTO \"token\"(user_id, token) VALUES($1, $2) RETURNING *;",
		userID, tokenVal)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (rep TokenRepositoryImpl) UpdateToken(userID uint32, tokenVal string) (*model.Token, error) {
	var token model.Token
	err := rep.db.Get(&token, "UPDATE token SET token = $1 WHERE user_id = $2 RETURNING *;",
		tokenVal, userID)
	if err != nil {
		return nil, err
	}
	return &token, nil
}
