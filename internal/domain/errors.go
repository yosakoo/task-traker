package domain

import "errors"

var (
	ErrUserNotFound            = errors.New("user doesn't exists")
	ErrUserAlreadyExists       = errors.New("user with such email already exists")
	ErrTokenExpired            = errors.New("token has expired")
	ErrTaskNotFound            = errors.New("task doesn't exists")
)