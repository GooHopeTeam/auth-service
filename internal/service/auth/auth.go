package auth

import (
	"github.com/goohopeteam/auth-service/internal/payload"
)

type AuthService interface {
	RegisterUser(pl *payload.RegistrationRequest) error
	LoginUser(pl *payload.LoginRequest) (*payload.TokenResponse, error)
	VerifyEmail(pl *payload.EmailVerificationRequest) (*payload.TokenResponse, error)
	VerifyToken(pl *payload.TokenVerificationRequest) error
}

var (
	EmailInUseErr            = payload.NewError("email_in_use")
	UserNotFoundErr          = payload.NewError("user_not_found")
	WrongCredentialsErr      = payload.NewError("wrong_credentials")
	InvalidTokenErr          = payload.NewError("invalid_token")
	WrongVerificationCodeErr = payload.NewError("wrong_verification_code")
)
