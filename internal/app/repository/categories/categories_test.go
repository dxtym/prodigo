package categories

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"prodigo/internal/app/models"
	"prodigo/pkg/db/postgres"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("can't run ddl", func(t *testing.T) {
		repo := new(postgres.MockPool)
		defer repo.AssertExpectations(t)

		repo.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, errors.New("cannot create table categories"))
		pool, err := New(repo)
		assert.NotNil(t, err)
		assert.Nil(t, pool)
	})
	t.Run("can run ddl", func(t *testing.T) {
		repo := new(postgres.MockPool)
		defer repo.AssertExpectations(t)

		repo.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.NewCommandTag("INSERT 1"), nil)
		pool, err := New(repo)
		assert.Nil(t, err)
		assert.NotNil(t, pool)
	})
}

func Test_repository_CreateCategory(t *testing.T) {
	t.Run("success on create category", func(t *testing.T) {
		mockRepo := new(MockRepo)

		c := models.Category{
			Name: "Test Category",
		}

		mockRepo.On("CreateCategory", mock.Anything, &c).Return(nil)
		err := mockRepo.CreateCategory(context.Background(), &c)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("error on create category", func(t *testing.T) {
		mockRepo := new(MockRepo)
		c := models.Category{}

		mockRepo.On("CreateCategory", mock.Anything, &c).Return(errors.New("db error"))
		err := mockRepo.CreateCategory(context.Background(), &c)
		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
		mockRepo.AssertExpectations(t)
	})
}
