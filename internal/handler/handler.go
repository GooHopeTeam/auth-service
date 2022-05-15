package handler

import (
	"errors"
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/goohopeteam/auth-service/internal/payload"
	"github.com/goohopeteam/auth-service/internal/service/auth"
)

func sendErrorResponse(ctx *gin.Context, err error) {
	if errResponse, ok := err.(payload.ErrorResponse); ok {
		ctx.JSON(http.StatusBadRequest, errResponse)
	} else {
		log.Printf("Enexpected error: %+v\n", err)
		ctx.JSON(http.StatusBadRequest, payload.NewUndefinedError())
	}
}

func extractPayload[Payload any](ctx *gin.Context) (*Payload, error) {
	details := map[string]string{
		"required": "required",
		"email":    "incorrect_email",
		"min":      "too_short",
		"max":      "too_long",
	}
	var payloadBody Payload
	payloadErr := payload.NewPayloadValidationError()
	if err := ctx.ShouldBindJSON(&payloadBody); err != nil {
		var vErrors validator.ValidationErrors
		if errors.As(err, &vErrors) {
			for _, vError := range vErrors {
				field, _ := reflect.TypeOf(payloadBody).FieldByName(vError.StructField())
				fieldName := field.Tag.Get("json")
				payloadErr.Details[fieldName] = details[vError.Tag()]
			}
		} else {
			return nil, err
		}
		return nil, payloadErr
	}
	return &payloadBody, nil
}

type Handler struct {
	authService auth.AuthService
}

func New(authService auth.AuthService) Handler {
	return Handler{authService: authService}
}

func (handler Handler) HandleRegistration(ctx *gin.Context) {
	pl, err := extractPayload[payload.RegistrationRequest](ctx)
	if err != nil {
		sendErrorResponse(ctx, err)
		return
	}
	err = handler.authService.RegisterUser(pl)
	if err != nil {
		sendErrorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, payload.EmptyResponse{})
}

func (handler Handler) HandleLogin(ctx *gin.Context) {
	var token *payload.TokenResponse
	pl, err := extractPayload[payload.LoginRequest](ctx)
	if err != nil {
		sendErrorResponse(ctx, err)
		return
	}
	token, err = handler.authService.LoginUser(pl)
	if err != nil {
		sendErrorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, token)
}

func (handler Handler) HandleTokenVerification(ctx *gin.Context) {
	pl, err := extractPayload[payload.TokenVerificationRequest](ctx)
	if err != nil {
		sendErrorResponse(ctx, err)
		return
	}
	err = handler.authService.VerifyToken(pl)
	if err != nil {
		sendErrorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, payload.EmptyResponse{})
}

func (handler Handler) HandleEmailVerification(ctx *gin.Context) {
	pl, err := extractPayload[payload.EmailVerificationRequest](ctx)
	if err != nil {
		sendErrorResponse(ctx, err)
		return
	}
	token, err := handler.authService.VerifyEmail(pl)
	if err != nil {
		sendErrorResponse(ctx, err)
	}
	ctx.JSON(http.StatusOK, token)
}
