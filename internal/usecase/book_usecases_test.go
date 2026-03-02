package usecase

import (
	"context"
	"errors"
	"task4/internal/domain"
	mockStorage "task4/internal/infrastructure/postgres/mock"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookCreate(t *testing.T) {
	mockStorage := new(mockStorage.BookStorageMock)

	testUC := &BookUseCase{
		bookRepo: mockStorage,
	}

	t.Run("success", func(t *testing.T) {
		expectedID := 123

		expectedBook := domain.Book{
			Author:          "test_author",
			Title:           "test_title",
			PublicationYear: 1999,
			Pages:           123,
			Genre:           "test_genre",
		}

		mockStorage.On("CreateWithLog", mock.Anything, expectedBook).Return(expectedID, nil).Once()

		id, err := testUC.Create(context.Background(), expectedBook)
		assert.NoError(t, err)
		assert.Equal(t, id, expectedID)
		mockStorage.AssertExpectations(t)
	})

	t.Run("failed", func(t *testing.T) {

		expErr := errors.New("definitely storage err")
		expId := 0
		mockStorage.On("CreateWithLog", mock.Anything, mock.Anything).Return(expId, expErr).Once()

		id, err := testUC.Create(context.Background(), domain.Book{})
		assert.Error(t, err)
		assert.Equal(t, id, expId)
		assert.Equal(t, err, expErr)
		mockStorage.AssertExpectations(t)
	})
}

func TestBookLoad(t *testing.T) {
	mockStorage := new(mockStorage.BookStorageMock)

	testUC := &BookUseCase{
		bookRepo: mockStorage,
	}

	t.Run("success", func(t *testing.T) {
		expectedID := 123

		expectedBook := domain.Book{
			ID:              123,
			Author:          "test_author",
			Title:           "test_title",
			PublicationYear: 1999,
			Pages:           123,
			Genre:           "test_genre",
		}

		mockStorage.On("Load", mock.Anything, expectedID).Return(expectedBook, nil).Once()

		book, err := testUC.Load(context.Background(), expectedID)
		assert.NoError(t, err)
		assert.Equal(t, book, expectedBook)
		mockStorage.AssertExpectations(t)
	})

	t.Run("failed", func(t *testing.T) {

		expErr := errors.New("definitely storage err")
		expId := 0
		expBook := domain.Book{}
		mockStorage.On("Load", mock.Anything, mock.Anything).Return(expBook, expErr).Once()

		book, err := testUC.Load(context.Background(), expId)
		assert.Error(t, err)
		assert.Equal(t, book, expBook)
		assert.Equal(t, err, expErr)
		mockStorage.AssertExpectations(t)
	})
}

func TestBookList(t *testing.T) {
	mockStorage := new(mockStorage.BookStorageMock)

	testUC := &BookUseCase{
		bookRepo: mockStorage,
	}

	t.Run("success", func(t *testing.T) {
		expBooks := []domain.Book{
			{
				ID:              1,
				Author:          "test_author1",
				Title:           "test_title1",
				PublicationYear: 1990,
				Pages:           120,
				Genre:           "test_genre1",
			},
			{
				ID:              2,
				Author:          "test_author2",
				Title:           "test_title2",
				PublicationYear: 1999,
				Pages:           123,
				Genre:           "test_genre1",
			},
			{
				ID:              3,
				Author:          "test_author3",
				Title:           "test_title3",
				PublicationYear: 1999,
				Pages:           120,
				Genre:           "test_genre2",
			},
		}

		expTotal := len(expBooks)

		mockStorage.On("List", mock.Anything, mock.Anything).Return(expBooks, expTotal, nil).Once()

		books, total, err := testUC.List(context.Background(), domain.BookFilter{})
		assert.NoError(t, err)
		assert.Equal(t, books, expBooks)
		assert.Equal(t, total, expTotal)
		mockStorage.AssertExpectations(t)
	})

	t.Run("failed", func(t *testing.T) {

		expErr := errors.New("definitely storage err")
		expTotal := 0

		mockStorage.On("List", mock.Anything, mock.Anything).Return([]domain.Book{}, expTotal, expErr).Once()

		books, total, err := testUC.List(context.Background(), domain.BookFilter{})
		assert.Error(t, err)
		assert.Equal(t, books, []domain.Book{})
		assert.Equal(t, total, expTotal)
		assert.Equal(t, err, expErr)
		mockStorage.AssertExpectations(t)
	})
}
