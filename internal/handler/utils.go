package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/goohopeteam/auth-service/internal/payload"
	"log"
	"net/http"
	"reflect"
)

func respondWithError(ctx *gin.Context, err error, code int) {
	if errResponse, ok := err.(payload.ErrorResponse); ok {
		ctx.JSON(code, errResponse)
	} else {
		log.Printf("Enexpected error: %+v\n", err)
		ctx.JSON(http.StatusInternalServerError, payload.NewUndefinedError())
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
		} else if err.Error() == "EOF" {
			return nil, payload.ErrorResponse{Err: "empty_body"}
		} else {
			return nil, err
		}
		return nil, payloadErr
	}
	return &payloadBody, nil
}
