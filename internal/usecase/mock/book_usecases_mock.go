package mock

import (
	"context"
	"task4/internal/domain"

	"github.com/stretchr/testify/mock"
)

type BookUseCaseMock struct {
	mock.Mock
}

func (m *BookUseCaseMock) List(ctx context.Context, filter domain.BookFilter) ([]domain.Book, int, error) {
	args := m.Called(ctx, filter)
	books, _ := args.Get(0).([]domain.Book)
	return books, args.Int(1), args.Error(2)
}

func (m *BookUseCaseMock) Load(ctx context.Context, id int) (domain.Book, error) {
	args := m.Called(ctx, id)
	book, _ := args.Get(0).(domain.Book)
	return book, args.Error(1)
}

func (m *BookUseCaseMock) Create(ctx context.Context, data domain.Book) (int, error) {
	args := m.Called(ctx, data)
	return args.Int(0), args.Error(1)
}

func (m *BookUseCaseMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *BookUseCaseMock) Update(ctx context.Context, id int, data domain.Book) error {
	args := m.Called(ctx, id, data)
	return args.Error(0)
}
