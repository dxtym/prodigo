package jwt

import "errors"

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("expired token")
	ErrInvalidSecretKey = errors.New("invalid secret key")
)
