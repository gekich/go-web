package errors

import "errors"

var (
	ErrEmailInUse         = errors.New("Email already in use")
	ErrInvalidCredentials = errors.New("Invalid credantials")
)
