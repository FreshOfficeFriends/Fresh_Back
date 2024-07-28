package domain

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

var ErrUserNotFound = errors.New("wrong email or password")
var TokenExpired = errors.New("token expired")
var EmptyTokenHeader = errors.New("empty token, add auth header")
var InvalidHeaderSignature = errors.New("invalid auth header, try 'Bearer <token>'")
var BadEmail = errors.New("empty 'email' field or invalid 'email' format")

func init() {
	validate = validator.New()
}

var validate *validator.Validate

type SignUp struct {
	FirstName  string `json:"first_name" validate:"required,lte=20,gte=2"`
	SecondName string `json:"last_name" validate:"required,lte=25,gte=2"`
	Email      string `json:"email" validate:"required,email"`
	Birthday   string `json:"birthday" validate:"required"`
	Password   string `json:"password" validate:"required,lte=64,gte=8,containsany=.!@#$%,containsany=1234567890"`
}

func (s SignUp) Validate() error {
	return validate.Struct(s)
}

type SignIn struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,lte=64,gte=8"`
}

type RefreshSession struct {
	ID        int64
	UserID    int
	Token     string
	ExpiresAt time.Time
}
