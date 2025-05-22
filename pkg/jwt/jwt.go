package jwt

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenMaker interface {
	CreateToken(int64, string, time.Duration) (string, error)
	VerifyToken(string) (*jwt.RegisteredClaims, error)
}

type tokenMaker struct {
	secretKey string
}

func New(secretKey string) (TokenMaker, error) {
	const minSecretKeyLength = 32
	if len(secretKey) < minSecretKeyLength {
		return nil, ErrInvalidSecretKey
	}

	return &tokenMaker{secretKey: secretKey}, nil
}

func (t *tokenMaker) CreateToken(userID int64, role string, duration time.Duration) (string, error) {
	payload := &jwt.RegisteredClaims{
		ID:        uuid.NewString(),
		Subject:   strconv.FormatInt(userID, 10),
		Audience:  jwt.ClaimStrings{role},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := token.SignedString([]byte(t.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func (t *tokenMaker) VerifyToken(token string) (*jwt.RegisteredClaims, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(t.secretKey), nil
	}

	claims := &jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !parsedToken.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
