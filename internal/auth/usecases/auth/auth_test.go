package auth_test

import (
	"context"
	"errors"
	"prodigo/internal/auth/dto"
	"prodigo/internal/auth/models"
	authRepository "prodigo/internal/auth/repository/auth"
	authService "prodigo/internal/auth/usecases/auth"
	"prodigo/pkg/jwt"
	"prodigo/pkg/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestService_Register(t *testing.T) {
	tests := []struct {
		name    string
		arg     dto.RegisterRequest
		wantErr error
	}{
		{
			name: "valid",
			arg: dto.RegisterRequest{
				Username: utils.GenerateRandomString(10),
				Password: utils.GenerateRandomString(10),
			},
			wantErr: nil,
		},
		{
			name: "invalid",
			arg: dto.RegisterRequest{
				Username: utils.GenerateRandomString(10),
				Password: utils.GenerateRandomString(10),
			},
			wantErr: errors.New("failed to register user"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := new(authRepository.MockRepository)
			defer repository.AssertExpectations(t)

			secretKey := utils.GenerateRandomString(32)
			maker, err := jwt.New(secretKey)
			require.NotNil(t, maker)
			require.NoError(t, err)

			repository.On("CreateUser",
				mock.Anything,
				mock.Anything,
			).Return(tt.wantErr).Once()

			service := authService.New(maker, repository)
			require.NotNil(t, service)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err = service.Register(ctx, tt.arg)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestService_Login(t *testing.T) {
	arg := dto.LoginRequest{
		Username: utils.GenerateRandomString(10),
		Password: utils.GenerateRandomString(10),
	}

	tests := []struct {
		name  string
		build func(*authRepository.MockRepository)
		check func(string, string, error)
	}{
		{
			name: "valid",
			build: func(repository *authRepository.MockRepository) {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(arg.Password), bcrypt.DefaultCost)
				require.NotEmpty(t, hashedPassword)
				require.NoError(t, err)

				repository.On("GetByUsername",
					mock.Anything,
					mock.Anything,
				).Return(&models.User{
					Username: arg.Username,
					Password: string(hashedPassword),
				}, nil).Once()

				repository.On("SaveToken",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(nil).Once()
			},
			check: func(accessToken, refreshToken string, err error) {
				assert.NotEmpty(t, accessToken)
				assert.NotEmpty(t, refreshToken)
				assert.NoError(t, err)
			},
		},
		{
			name: "not found",
			build: func(repository *authRepository.MockRepository) {
				repository.On("GetByUsername",
					mock.Anything,
					mock.Anything,
				).Return(&models.User{}, authRepository.ErrUserNotFound).Once()
			},
			check: func(accessToken, refreshToken string, err error) {
				assert.Empty(t, accessToken)
				assert.Empty(t, refreshToken)
				assert.ErrorContains(t, err, authRepository.ErrUserNotFound.Error())
			},
		},
		{
			name: "invalid credentials",
			build: func(repository *authRepository.MockRepository) {
				wrongPassword := utils.GenerateRandomString(10)
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(wrongPassword), bcrypt.DefaultCost)
				require.NotEmpty(t, hashedPassword)
				require.NoError(t, err)

				repository.On("GetByUsername",
					mock.Anything,
					mock.Anything,
				).Return(&models.User{
					Username: arg.Username,
					Password: string(hashedPassword),
				}, nil).Once()
			},
			check: func(accessToken, refreshToken string, err error) {
				assert.Empty(t, accessToken)
				assert.Empty(t, refreshToken)
				assert.ErrorContains(t, err, authService.ErrInvalidCredentials.Error())
			},
		},
		{
			name: "invalid token",
			build: func(repository *authRepository.MockRepository) {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(arg.Password), bcrypt.DefaultCost)
				require.NotEmpty(t, hashedPassword)
				require.NoError(t, err)

				repository.On("GetByUsername",
					mock.Anything,
					mock.Anything,
				).Return(&models.User{
					Username: arg.Username,
					Password: string(hashedPassword),
				}, nil).Once()

				repository.On("SaveToken",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(errors.New("failed to save token")).Once()
			},
			check: func(accessToken, refreshToken string, err error) {
				assert.Empty(t, accessToken)
				assert.Empty(t, refreshToken)
				assert.ErrorContains(t, err, errors.New("failed to save token").Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := new(authRepository.MockRepository)
			defer repository.AssertExpectations(t)

			secretKey := utils.GenerateRandomString(32)
			maker, err := jwt.New(secretKey)
			require.NotNil(t, maker)
			require.NoError(t, err)

			tt.build(repository)

			service := authService.New(maker, repository)
			require.NotNil(t, service)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			accessToken, refreshToken, err := service.Login(ctx, arg)
			tt.check(accessToken, refreshToken, err)
		})
	}
}

func TestService_Refresh(t *testing.T) {
	tests := []struct {
		name  string
		build func(*authRepository.MockRepository, string)
		check func(string, error)
	}{
		{
			name: "valid",
			build: func(repository *authRepository.MockRepository, refreshToken string) {
				repository.On("GetToken",
					mock.Anything,
					mock.Anything,
				).Return(refreshToken, nil).Once()
			},
			check: func(accessToken string, err error) {
				assert.NotEmpty(t, accessToken)
				assert.NoError(t, err)
			},
		},
		{
			name: "not found",
			build: func(repository *authRepository.MockRepository, refreshToken string) {
				repository.On("GetToken",
					mock.Anything,
					mock.Anything,
				).Return("", authService.ErrTokenNotFound).Once()
			},
			check: func(accessToken string, err error) {
				assert.Empty(t, accessToken)
				assert.ErrorContains(t, err, authService.ErrTokenNotFound.Error())
			},
		},
		{
			name: "invalid token",
			build: func(repository *authRepository.MockRepository, refreshToken string) {
				wrongToken := utils.GenerateRandomString(10)
				repository.On("GetToken",
					mock.Anything,
					mock.Anything,
				).Return(wrongToken, nil).Once()
			},
			check: func(accessToken string, err error) {
				assert.Empty(t, accessToken)
				assert.ErrorContains(t, err, authService.ErrInvalidToken.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := new(authRepository.MockRepository)
			defer repository.AssertExpectations(t)

			secretKey := utils.GenerateRandomString(32)
			maker, err := jwt.New(secretKey)
			require.NotNil(t, maker)
			require.NoError(t, err)

			refreshToken, err := maker.CreateToken(
				utils.GenerateRandomInt(10),
				utils.GenerateRandomString(10),
				time.Minute,
			)
			require.NotNil(t, refreshToken)
			require.NoError(t, err)

			arg := dto.RefreshRequest{RefreshToken: refreshToken}

			tt.build(repository, refreshToken)

			service := authService.New(maker, repository)
			require.NotNil(t, service)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			accessToken, err := service.Refresh(ctx, arg)
			tt.check(accessToken, err)
		})
	}
}
