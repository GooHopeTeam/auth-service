package auth

import (
	"github.com/goohopeteam/auth-service/internal/model"
	"github.com/goohopeteam/auth-service/internal/repository"
	"github.com/goohopeteam/auth-service/internal/service/verifier"
	"time"
)

const hashSalt = "abc"
const verificationCode = "123"

type UserRepositoryMock struct {
	repository.UserRepository
	db []model.User
}

func (rep *UserRepositoryMock) FindByEmail(email string) (*model.User, error) {
	return findInSlice[model.User](rep.db, func(user *model.User) bool {
		return user.Email == email
	}), nil
}

func (rep *UserRepositoryMock) Insert(email, hashedPassword string) (*model.User, error) {
	user := model.User{ID: uint32(len(rep.db)) + 1, Email: email, HashedPassword: hashedPassword, CreatedAt: time.Now()}
	rep.db = append(rep.db, user)
	return &user, nil
}

type TokenRepositoryMock struct {
	repository.TokenRepository
	db []model.Token
}

func (rep *TokenRepositoryMock) Find(userId uint32) (*model.Token, error) {
	return findInSlice[model.Token](rep.db, func(token *model.Token) bool {
		return token.UserID == userId
	}), nil
}

func (rep *TokenRepositoryMock) Insert(user *model.User, tokenVal string) (*model.Token, error) {
	token := model.Token{UserID: user.ID, Value: tokenVal}
	rep.db = append(rep.db, token)
	return &token, nil
}

func initMockRepositories() (UserRepositoryMock, TokenRepositoryMock) {
	userDb := []model.User{
		{ID: 1, Email: "user1@gmail.com", HashedPassword: makeHash("123456", hashSalt), CreatedAt: time.Now()},
		{ID: 2, Email: "user2@gmail.com", HashedPassword: makeHash("qwerty", hashSalt), CreatedAt: time.Now()},
	}
	tokenDb := []model.Token{
		{UserID: 1, Value: func() string { token, _ := generateToken(&userDb[0]); return token }()},
		{UserID: 2, Value: func() string { token, _ := generateToken(&userDb[1]); return token }()},
	}
	return UserRepositoryMock{db: userDb}, TokenRepositoryMock{db: tokenDb}
}

type VerifierMock struct {
	verifier.Verifier
	storage map[string]map[string]string
}

func newVerifierMock() *VerifierMock {
	return &VerifierMock{storage: make(map[string]map[string]string)}
}

func (v VerifierMock) Send(email string, data map[string]string) error {
	v.storage[email] = data
	return nil
}

func (v VerifierMock) Check(email, code string) (map[string]string, error) {
	if code == verificationCode {
		return v.storage[email], nil
	}

	return nil, nil
}
