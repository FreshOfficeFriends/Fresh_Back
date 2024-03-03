package rest

import "github.com/go-playground/validator/v10"

type StatusCode string

func init() {
	validate = validator.New()
}

var validate *validator.Validate

const (
	success StatusCode = "success"
	fail    StatusCode = "fail"
)

type Response struct {
	Status StatusCode    `json:"statusCode"`
	Data   interface{}   `json:"data,omitempty"`
	Error  *ErrorDetails `json:"error,omitempty"`
}

// ErrorDetails содержит информацию об ошибке
type ErrorDetails struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type Email struct {
	Email string `json:"email" validate:"required,email"`
}

func (e Email) Validate() error {
	return validate.Struct(e)
}
