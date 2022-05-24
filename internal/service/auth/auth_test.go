package auth

import (
	"github.com/goohopeteam/auth-service/internal/payload"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getAuthService() AuthServiceImpl {
	userRep, tokenRep := initMockRepositories()
	return AuthServiceImpl{userRep: &userRep, tokenRep: &tokenRep, verifier: newVerifierMock(), globalSalt: hashSalt}
}

func TestAuthServiceImpl_RegisterUser_Error_EmailInUse(t *testing.T) {
	assert := assert.New(t)
	regRequest := &payload.RegistrationRequest{
		Email:    "user1@gmail.com",
		Password: "abc123",
	}
	err := getAuthService().RegisterUser(regRequest)
	assert.ErrorIs(err, EmailInUseErr)
}

func TestAuthServiceImpl_RegisterUser_Success(t *testing.T) {
	assert := assert.New(t)
	regRequest := &payload.RegistrationRequest{
		Email:    "user3@gmail.com",
		Password: "abc123",
	}
	authService := getAuthService()
	err := authService.RegisterUser(regRequest)
	assert.NoError(err)
	data, err := authService.verifier.Check(regRequest.Email, verificationCode)
	assert.NotNil(data)
	assert.NoError(err)
	assert.Equal(makeHash(regRequest.Password, hashSalt), data["hashedPassword"])
}

func TestAuthServiceImpl_VerifyEmail_Error_WrongCode(t *testing.T) {
	assert := assert.New(t)
	verRequest := &payload.EmailVerificationRequest{
		Email: "user3@gmail.com",
		Code:  "111",
	}
	authService := getAuthService()
	err := authService.verifier.Send(verRequest.Email, map[string]string{"hashedPassword": "qwerty"})
	assert.NoError(err)
	data, err := authService.VerifyEmail(verRequest)
	assert.NotNil(err)
	assert.Nil(data)
	assert.ErrorIs(err, WrongVerificationCodeErr)
}

func TestAuthServiceImpl_VerifyEmail_Error_WrongEmail(t *testing.T) {
	assert := assert.New(t)
	verRequest := &payload.EmailVerificationRequest{
		Email: "user4@gmail.com",
		Code:  verificationCode,
	}
	authService := getAuthService()
	data, err := authService.VerifyEmail(verRequest)
	assert.NotNil(err)
	assert.Nil(data)
	assert.ErrorIs(err, WrongVerificationCodeErr)
}

func TestAuthServiceImpl_VerifyEmail_Success(t *testing.T) {
	assert := assert.New(t)
	verRequest := &payload.EmailVerificationRequest{
		Email: "user3@gmail.com",
		Code:  verificationCode,
	}
	authService := getAuthService()
	err := authService.verifier.Send(verRequest.Email, map[string]string{"hashedPassword": "qwerty"})
	assert.NoError(err)
	data, err := authService.VerifyEmail(verRequest)
	assert.NotNil(data)
	assert.NoError(err)
	user, err := authService.userRep.FindByEmail("user3@gmail.com")
	assert.NoError(err)
	assert.NotNil(user)
	assert.Equal(uint32(3), user.Id)
}

func TestAuthServiceImpl_LoginUser_Error_WrongCredentials_WrongEmail(t *testing.T) {
	assert := assert.New(t)
	loginRequest := &payload.LoginRequest{
		Email:    "user3@gmail.com",
		Password: "123456",
	}
	authService := getAuthService()
	_, err := authService.LoginUser(loginRequest)
	assert.ErrorIs(err, WrongCredentialsErr)
}

func TestAuthServiceImpl_LoginUser_Error_WrongCredentials_WrongPassword(t *testing.T) {
	assert := assert.New(t)
	loginRequest := &payload.LoginRequest{
		Email:    "user1@gmail.com",
		Password: "wrong_password",
	}
	authService := getAuthService()
	_, err := authService.LoginUser(loginRequest)
	assert.ErrorIs(err, WrongCredentialsErr)
}

func TestAuthServiceImpl_LoginUser_Success(t *testing.T) {
	assert := assert.New(t)
	loginRequest := &payload.LoginRequest{
		Email:    "user1@gmail.com",
		Password: "123456",
	}
	authService := getAuthService()
	_, err := authService.LoginUser(loginRequest)
	assert.NoError(err)
}

func TestAuthServiceImpl_VerifyToken_WrongUser(t *testing.T) {
	assert := assert.New(t)
	verRequest := &payload.TokenVerificationRequest{
		UserId: 3,
		Token:  "123",
	}
	authService := getAuthService()
	err := authService.VerifyToken(verRequest)
	assert.ErrorIs(err, InvalidTokenErr)
}

func TestAuthServiceImpl_VerifyToken_WrongToken(t *testing.T) {
	assert := assert.New(t)
	verRequest := &payload.TokenVerificationRequest{
		UserId: 1,
		Token:  "123",
	}
	authService := getAuthService()
	err := authService.VerifyToken(verRequest)
	assert.ErrorIs(err, InvalidTokenErr)
}

func TestAuthServiceImpl_VerifyToken_Success(t *testing.T) {
	assert := assert.New(t)
	authService := getAuthService()
	verRequest := &payload.TokenVerificationRequest{
		UserId: 1,
		Token:  func() string { token, _ := authService.tokenRep.Find(1); return token.Value }(),
	}
	err := authService.VerifyToken(verRequest)
	assert.NoError(err)
}
