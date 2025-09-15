package user

import (
	"fmt"
	"strings"
)

type UserExistsError struct {
	level string
	Err   error
}
type AuthorizationError struct {
	level string
	Err   error
}

type UserNotFoundError struct {
	level string
	Err   error
}

func (le *UserExistsError) Error() string {
	return fmt.Sprintf("[%s] %v", le.level, le.Err)
}

func NewUserExistsError(label string, err error) error {
	return &UserExistsError{
		level: strings.ToUpper(label),
		Err:   err,
	}
}

func (le *AuthorizationError) Error() string {
	return fmt.Sprintf("[%s] %v", le.level, le.Err)
}

func NewAuthorizationError(label string, err error) error {
	return &AuthorizationError{
		level: strings.ToUpper(label),
		Err:   err,
	}
}

func (le *UserNotFoundError) Error() string {
	return fmt.Sprintf("[%s] %v", le.level, le.Err)
}

func NewUserNotFoundError(label string, err error) error {
	return &UserNotFoundError{
		level: strings.ToUpper(label),
		Err:   err,
	}
}
