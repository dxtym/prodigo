package jwt_test

import (
	maker "prodigo/pkg/jwt"
	"prodigo/pkg/utils"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		secretKey string
		wantErr   error
	}{
		{name: "valid secret", secretKey: utils.GenerateRandomString(32), wantErr: nil},
		{name: "invalid secret", secretKey: utils.GenerateRandomString(10), wantErr: maker.ErrInvalidSecretKey},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := maker.New(tt.secretKey)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCreateAndVerifyToken(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		secretKey := utils.GenerateRandomString(32)
		m, err := maker.New(secretKey)
		require.NotNil(t, m)
		require.NoError(t, err)

		userID := utils.GenerateRandomInt(100)
		role := utils.GenerateRandomString(10)
		duration := time.Minute

		token, err := m.CreateToken(userID, role, duration)
		require.NotEmpty(t, token)
		require.NoError(t, err)

		payload, err := m.VerifyToken(token)
		require.NotNil(t, payload)
		require.NoError(t, err)

		assert.NotZero(t, payload.ID)
		assert.Equal(t, strconv.FormatInt(userID, 10), payload.Subject)

		require.NotEmpty(t, payload.Audience)
		assert.Equal(t, role, payload.Audience[0])

		assert.WithinDuration(t, time.Now(), payload.IssuedAt.Time, time.Second)
		assert.WithinDuration(t, time.Now().Add(duration), payload.ExpiresAt.Time, time.Second)
	})

	t.Run("expired token", func(t *testing.T) {
		secretKey := utils.GenerateRandomString(32)
		m, err := maker.New(secretKey)
		require.NotNil(t, m)
		require.NoError(t, err)

		userID := utils.GenerateRandomInt(100)
		role := utils.GenerateRandomString(10)
		duration := -time.Minute

		token, err := m.CreateToken(userID, role, duration)
		require.NotEmpty(t, token)
		require.NoError(t, err)

		payload, err := m.VerifyToken(token)
		assert.Nil(t, payload)
		assert.ErrorIs(t, err, maker.ErrExpiredToken)
	})

	t.Run("invalid token", func(t *testing.T) {
		inPayload := &jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Subject:   utils.GenerateRandomString(100),
			Audience:  jwt.ClaimStrings{utils.GenerateRandomString(10)},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodNone, inPayload)
		require.NotNil(t, token)

		signedToken, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
		require.NotEmpty(t, signedToken)
		require.NoError(t, err)

		secretKey := utils.GenerateRandomString(32)
		m, err := maker.New(secretKey)
		require.NotNil(t, m)
		require.NoError(t, err)

		outPaylod, err := m.VerifyToken(signedToken)
		assert.Nil(t, outPaylod)
		assert.ErrorIs(t, err, maker.ErrInvalidToken)
	})
}
