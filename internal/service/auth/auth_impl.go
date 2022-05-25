package auth

import (
	"github.com/goohopeteam/auth-service/internal/payload"
	"github.com/goohopeteam/auth-service/internal/repository"
	"github.com/goohopeteam/auth-service/internal/service/verifier"
)

type AuthServiceImpl struct {
	AuthService
	userRep    repository.UserRepository
	tokenRep   repository.TokenRepository
	verifier   verifier.Verifier
	globalSalt string
}

func New(userRep repository.UserRepository, tokenRep repository.TokenRepository, verifier verifier.Verifier, globalSalt string) AuthService {
	return AuthServiceImpl{userRep: userRep, tokenRep: tokenRep,
		globalSalt: globalSalt, verifier: verifier}
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
	err = s.verifier.Send(pl.Email, map[string]string{"hashedPassword": hash})
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

	token, err := s.tokenRep.Find(user.ID)
	return &payload.TokenResponse{UserID: user.ID, Value: token.Value}, err
}

func (s AuthServiceImpl) VerifyEmail(pl *payload.EmailVerificationRequest) (*payload.TokenResponse, error) {
	data, err := s.verifier.Check(pl.Email, pl.Code)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, WrongVerificationCodeErr
	}

	user, err := s.userRep.Insert(pl.Email, data["hashedPassword"])
	if err != nil {
		return nil, err
	}

	tokenVal, err := generateToken(user)
	if err != nil {
		return nil, err
	}

	token, err := s.tokenRep.Insert(user.ID, tokenVal)
	if err != nil {
		return nil, err
	}

	return &payload.TokenResponse{UserID: user.ID, Value: token.Value}, nil
}

func (s AuthServiceImpl) VerifyToken(pl *payload.TokenVerificationRequest) error {
	token, err := s.tokenRep.Find(pl.UserID)
	if err != nil {
		return err
	}

	if token == nil || token.Value != pl.Token {
		return InvalidTokenErr
	}

	return nil
}

func (s AuthServiceImpl) ChangePassword(pl *payload.ChangePasswordRequest) (*payload.TokenResponse, error) {
	user, err := s.userRep.FindByEmail(pl.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, WrongCredentialsErr
	}

	oldHashedPassword := makeHash(pl.OldPassword, s.globalSalt)
	if user.HashedPassword != oldHashedPassword {
		return nil, WrongCredentialsErr
	}

	newHashedPassword := makeHash(pl.NewPassword, s.globalSalt)

	user, err = s.userRep.UpdatePassword(user.ID, newHashedPassword)
	if err != nil {
		return nil, err
	}

	tokenVal, err := generateToken(user)
	if err != nil {
		return nil, err
	}

	token, err := s.tokenRep.UpdateToken(user.ID, tokenVal)
	if err != nil {
		return nil, err
	}

	return &payload.TokenResponse{UserID: user.ID, Value: token.Value}, nil
}
