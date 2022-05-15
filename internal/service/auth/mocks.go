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
	user := model.User{Id: uint32(len(rep.db)) + 1, Email: email, HashedPassword: hashedPassword, CreatedAt: time.Now()}
	rep.db = append(rep.db, user)
	return &user, nil
}

type TokenRepositoryMock struct {
	repository.TokenRepository
	db []model.Token
}

func (rep *TokenRepositoryMock) Find(userId uint32) (*model.Token, error) {
	return findInSlice[model.Token](rep.db, func(token *model.Token) bool {
		return token.UserId == userId
	}), nil
}

func (rep *TokenRepositoryMock) Insert(user *model.User, tokenVal string) (*model.Token, error) {
	token := model.Token{UserId: user.Id, Value: tokenVal}
	rep.db = append(rep.db, token)
	return &token, nil
}

func initMockRepositories() (UserRepositoryMock, TokenRepositoryMock) {
	userDb := []model.User{
		{1, "user1@gmail.com", makeHash("123456", hashSalt), time.Now()},
		{2, "user2@gmail.com", makeHash("qwerty", hashSalt), time.Now()},
	}
	tokenDb := []model.Token{
		{1, func() string { token, _ := generateToken(&userDb[0]); return token }()},
		{2, func() string { token, _ := generateToken(&userDb[1]); return token }()},
	}
	return UserRepositoryMock{db: userDb}, TokenRepositoryMock{db: tokenDb}
}

type VerifierMock struct {
	verifier.Verifier[verificationData]
	storage map[string]*verificationData
}

func (v VerifierMock) Send(email string, data *verificationData) error {
	v.storage[email] = data
	return nil
}

func (v VerifierMock) Check(email, code string) (*verificationData, error) {
	if code == verificationCode {
		return v.storage[email], nil
	}
	return nil, nil
}

func getAuthService() AuthServiceImpl {
	userRep, tokenRep := initMockRepositories()
	v := VerifierMock{storage: make(map[string]*verificationData)}
	return AuthServiceImpl{userRep: &userRep, tokenRep: &tokenRep, verifier: v, globalSalt: hashSalt}
}
