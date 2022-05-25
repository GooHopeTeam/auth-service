package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/goohopeteam/auth-service/internal/payload"
	"github.com/goohopeteam/auth-service/internal/service/auth"
	"net/http"
)

type Handler struct {
	authService auth.AuthService
}

func New(authService auth.AuthService) Handler {
	return Handler{authService: authService}
}

func (handler Handler) HandleRegistration(ctx *gin.Context) {
	pl, err := extractPayload[payload.RegistrationRequest](ctx)
	if err != nil {
		respondWithError(ctx, err, http.StatusBadRequest)
		return
	}
	err = handler.authService.RegisterUser(pl)
	if err != nil {
		respondWithError(ctx, err, http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, payload.EmptyResponse{})
}

func (handler Handler) HandleLogin(ctx *gin.Context) {
	var token *payload.TokenResponse
	pl, err := extractPayload[payload.LoginRequest](ctx)
	if err != nil {
		respondWithError(ctx, err, http.StatusBadRequest)
		return
	}
	token, err = handler.authService.LoginUser(pl)
	if err != nil {
		respondWithError(ctx, err, http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, token)
}

func (handler Handler) HandleTokenVerification(ctx *gin.Context) {
	pl, err := extractPayload[payload.TokenVerificationRequest](ctx)
	if err != nil {
		respondWithError(ctx, err, http.StatusBadRequest)
		return
	}
	err = handler.authService.VerifyToken(pl)
	if err != nil {
		respondWithError(ctx, err, http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, payload.EmptyResponse{})
}

func (handler Handler) HandleEmailVerification(ctx *gin.Context) {
	pl, err := extractPayload[payload.EmailVerificationRequest](ctx)
	if err != nil {
		respondWithError(ctx, err, http.StatusBadRequest)
		return
	}
	token, err := handler.authService.VerifyEmail(pl)
	if err != nil {
		respondWithError(ctx, err, http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, token)
}

func (handler Handler) HandlePasswordChange(ctx *gin.Context) {
	pl, err := extractPayload[payload.ChangePasswordRequest](ctx)
	if err != nil {
		respondWithError(ctx, err, http.StatusBadRequest)
		return
	}
	token, err := handler.authService.ChangePassword(pl)
	if err != nil {
		respondWithError(ctx, err, http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, token)
}
