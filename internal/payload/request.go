package payload

type RegistrationRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=80"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=80"`
}

type TokenVerificationRequest struct {
	UserID uint32 `json:"user_id" binding:"required"`
	Token  string `json:"token" binding:"required"`
}

type EmailVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}
