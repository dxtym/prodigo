package auth

import (
	"context"
	"errors"
	"fmt"
	"prodigo/internal/auth/dto"
	"prodigo/internal/auth/models"
	"prodigo/internal/auth/repository/auth"
	"prodigo/pkg/jwt"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	accessDuration  = 15 * time.Minute
	refreshDuration = 24 * time.Hour
)

type Service interface {
	Register(context.Context, dto.RegisterRequest) error
	Login(context.Context, dto.LoginRequest) (string, string, error)
	Refresh(context.Context, dto.RefreshRequest) (string, error)
}

type service struct {
	maker      jwt.TokenMaker
	repository auth.Repository
}

func New(maker jwt.TokenMaker, repository auth.Repository) Service {
	return &service{maker: maker, repository: repository}
}

func (s *service) Register(ctx context.Context, req dto.RegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	if err := s.repository.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	return nil
}

func (s *service) Login(ctx context.Context, req dto.LoginRequest) (access, refresh string, err error) {
	user, err := s.repository.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			return "", "", ErrUserNotFound
		}
		return "", "", fmt.Errorf("failed to get by username: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	accessToken, err := s.maker.CreateToken(user.ID, user.Role, accessDuration)
	if err != nil {
		return "", "", fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := s.maker.CreateToken(user.ID, user.Role, refreshDuration)
	if err != nil {
		return "", "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	if err := s.repository.SaveToken(ctx, user.ID, refreshToken, refreshDuration); err != nil {
		return "", "", fmt.Errorf("failed to save refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *service) Refresh(ctx context.Context, req dto.RefreshRequest) (string, error) {
	payload, err := s.maker.VerifyToken(req.RefreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrExpiredToken) {
			return "", ErrExpiredToken
		}
		return "", fmt.Errorf("failed to verify refresh token: %w", err)
	}

	userRole := payload.Audience[0]
	userID, err := strconv.ParseInt(payload.Subject, 10, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse user id from token: %w", err)
	}

	refreshToken, err := s.repository.GetToken(ctx, userID)
	if err != nil {
		if errors.Is(err, auth.ErrTokenNotFound) {
			return "", ErrTokenNotFound
		}
		return "", fmt.Errorf("failed to get refresh token: %w", err)
	}

	if refreshToken != req.RefreshToken {
		return "", ErrInvalidToken
	}

	accessToken, err := s.maker.CreateToken(userID, userRole, accessDuration)
	if err != nil {
		return "", fmt.Errorf("failed to create access token: %w", err)
	}

	return accessToken, nil
}
