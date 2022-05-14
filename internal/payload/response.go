package payload

type EmptyResponse struct {
}

type TokenResponse struct {
	UserId string `json:"user_id"`
	Value  string `json:"token"`
}

type ErrorResponse struct {
	error
	Err     string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

func (e ErrorResponse) Error() string {
	return e.Err
}

func NewError(err string) ErrorResponse {
	return ErrorResponse{Err: err}
}

func NewPayloadValidationError() ErrorResponse {
	return ErrorResponse{Err: "payload_validation", Details: make(map[string]string)}
}

func NewUndefinedError() ErrorResponse {
	return ErrorResponse{Err: "undefined"}
}