package auth

import (
	"github.com/goohopeteam/auth-service/internal/service/verifier"
	"strconv"

	"github.com/goohopeteam/auth-service/internal/payload"
	"github.com/goohopeteam/auth-service/internal/repository"
)

type verificationData struct {
	hashedPassword string
}

type AuthServiceImpl struct {
	AuthService
	userRep    repository.UserRepository
	tokenRep   repository.TokenRepository
	verifier   verifier.Verifier[verificationData]
	globalSalt string
}

func New(userRep repository.UserRepository, tokenRep repository.TokenRepository, globalSalt string) AuthService {
	return AuthServiceImpl{userRep: userRep, tokenRep: tokenRep,
		globalSalt: globalSalt, verifier: verifier.EmailVerifier[verificationData]{}}
}

func (s AuthServiceImpl) RegisterUser(pl *payload.RegistrationRequest) error {
	user, err := s.userRep.FindByEmail(pl.Email)
	if err != nil {
		return err
	}
	if user != nil {
		return EmailInUseErr
	}
	hash := makeHash(pl.Password, s.globalSalt)
	err = s.verifier.Send(pl.Email, &verificationData{hashedPassword: hash})
	return err
}

func (s AuthServiceImpl) LoginUser(pl *payload.LoginRequest) (*payload.TokenResponse, error) {
	user, err := s.userRep.FindByEmail(pl.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, WrongCredentialsErr
	}
	hash := makeHash(pl.Password, s.globalSalt)
	if user.HashedPassword != hash {
		return nil, WrongCredentialsErr
	}

	token, err := s.tokenRep.Find(user.Id)
	return &payload.TokenResponse{UserId: strconv.Itoa(int(user.Id)), Value: token.Value}, err
}

func (s AuthServiceImpl) VerifyEmail(pl *payload.EmailVerificationRequest) (*payload.TokenResponse, error) {
	data, err := s.verifier.Check(pl.Email, pl.Code)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, WrongVerificationCodeErr
	}
	user, err := s.userRep.Insert(pl.Email, (*data).hashedPassword)
	if err != nil {
		return nil, err
	}
	tokenVal, err := generateToken(user)
	if err != nil {
		return nil, err
	}
	token, err := s.tokenRep.Insert(user, tokenVal)
	if err != nil {
		return nil, err
	}
	return &payload.TokenResponse{Value: token.Value}, nil
}

func (s AuthServiceImpl) VerifyToken(pl *payload.TokenVerificationRequest) error {
	token, err := s.tokenRep.Find(pl.UserId)
	if err != nil {
		return err
	}
	if token == nil {
		return InvalidTokenErr
	}
	if token.Value != pl.Token {
		return InvalidTokenErr
	}
	return nil
}
