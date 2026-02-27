package usecase

import (
	"context"
	"log"
	"task4/internal/domain"
)

type BookUseCase struct {
	bookRepo BookRepository
}

type BookRepository interface {
	Create(ctx context.Context, data domain.Book) (id int, err error)
	List(ctx context.Context, f domain.BookFilter) ([]domain.Book, int, error)
	Load(ctx context.Context, id int) (domain.Book, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, data domain.Book) error
}

func NewBookUseCase(bookRepo BookRepository) BookUseCase {
	return BookUseCase{
		bookRepo: bookRepo,
	}
}

func (uc *BookUseCase) List(ctx context.Context, filter domain.BookFilter) ([]domain.Book, int, error) {
	if filter.YearTo > 0 && filter.YearFrom > filter.YearTo {
		return nil, 0, domain.ErrInvalidFilterRange
	}
	books, total, err := uc.bookRepo.List(ctx, filter)
	if err != nil {
		log.Println("book usecase list err:", err)
		return books, total, err
	}
	return books, total, nil
}

func (uc *BookUseCase) Load(ctx context.Context, id int) (domain.Book, error) {
	book, err := uc.bookRepo.Load(ctx, id)
	if err != nil {
		log.Println("book usecase load err:", err)
		return book, err
	}
	return book, nil
}

func (uc *BookUseCase) Create(ctx context.Context, data domain.Book) (int, error) {
	log.Println("new data:", data)
	id, err := uc.bookRepo.Create(ctx, data)
	if err != nil {
		log.Println("book usecase create err:", err)
		return id, err
	}
	return id, nil
}

func (uc *BookUseCase) Delete(ctx context.Context, id int) error {
	err := uc.bookRepo.Delete(ctx, id)
	if err != nil {
		log.Println("book usecase delete err:", err)
		return err
	}
	return nil
}

func (uc *BookUseCase) Update(ctx context.Context, id int, data domain.Book) error {
	err := uc.bookRepo.Update(ctx, id, data)
	if err != nil {
		log.Println("book usecase update err:", err)
		return err
	}
	return nil
}
