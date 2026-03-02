package postgres

import (
	"context"
	"errors"
	"fmt"
	"task4/internal/domain"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookStorage struct {
	db *pgxpool.Pool
}

func NewBookStorage(pool *pgxpool.Pool) *BookStorage {
	return &BookStorage{
		db: pool,
	}
}

func (s *BookStorage) CreateWithLog(ctx context.Context, data domain.Book) (id int, err error) {

	// begin transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return id, fmt.Errorf("db create: failed to create transaction: %w", err)
	}

	// rollback transaction
	defer tx.Rollback(ctx)

	// create new book
	bookInsert := squirrel.Insert("books").
		Columns("author", "title", "publication_year", "pages", "genre").
		Values(data.Author, data.Title, data.PublicationYear, data.Pages, data.Genre).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar)

	bookInsertSql, args, err := bookInsert.ToSql()
	if err != nil {
		return id, fmt.Errorf("db create: failed to generate SQL (book insert): %w", err)
	}

	row := tx.QueryRow(ctx, bookInsertSql, args...)
	err = row.Scan(&id)
	if err != nil {
		return id, fmt.Errorf("db create: insert book error: %w", err)
	}

	// create new log recording
	createLog := squirrel.Insert("books_log").
		Columns("book_id", "action").
		Values(id, "create").
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar)

	createLogSql, args, err := createLog.ToSql()
	if err != nil {
		return id, fmt.Errorf("db create: failed to generate SQL (book_log insert): %w", err)
	}

	row = tx.QueryRow(ctx, createLogSql, args...)
	err = row.Scan(&id)
	if err != nil {
		return id, fmt.Errorf("db create: insert book_log error: %w", err)
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return id, fmt.Errorf("db create: books create commit tx error: %w", err)
	}

	return id, nil
}

func (s *BookStorage) List(ctx context.Context, f domain.BookFilter) ([]domain.Book, int, error) {

	conditions := squirrel.And{}
	if f.Author != "" {
		conditions = append(conditions, squirrel.ILike{"author": "%" + f.Author + "%"})
	}
	if f.Title != "" {
		conditions = append(conditions, squirrel.ILike{"title": "%" + f.Title + "%"})
	}
	if f.Genre != "" {
		conditions = append(conditions, squirrel.ILike{"genre": "%" + f.Genre + "%"})
	}
	if f.YearFrom > 0 {
		conditions = append(conditions, squirrel.GtOrEq{"publication_year": f.YearFrom})
	}
	if f.YearTo > 0 {
		conditions = append(conditions, squirrel.LtOrEq{"publication_year": f.YearTo})
	}
	if f.PagesFrom > 0 {
		conditions = append(conditions, squirrel.GtOrEq{"pages": f.PagesFrom})
	}
	if f.PagesTo > 0 {
		conditions = append(conditions, squirrel.LtOrEq{"pages": f.PagesTo})
	}

	countQuery, args, err := squirrel.Select("COUNT(*)").
		From("books").
		Where(conditions).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, 0, fmt.Errorf("build count sql: %w", err)
	}

	var total int
	if err := s.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("exec count sql: %w", err)
	}

	if total == 0 {
		return []domain.Book{}, 0, nil
	}

	dataQuery, args, err := squirrel.Select("id", "author", "title", "publication_year", "pages", "genre").
		From("books").
		Where(conditions).
		OrderBy("id").
		Limit(f.Limit).
		Offset(f.Offset).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, 0, fmt.Errorf("build data sql: %w", err)
	}

	rows, err := s.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("exec data query: %w", err)
	}
	defer rows.Close()

	var books []domain.Book
	for rows.Next() {
		var b domain.Book
		err := rows.Scan(&b.ID, &b.Author, &b.Title, &b.PublicationYear, &b.Pages, &b.Genre)
		if err != nil {
			return nil, 0, fmt.Errorf("scan book: %w", err)
		}
		books = append(books, b)
	}

	return books, total, nil
}

func (s *BookStorage) Load(ctx context.Context, id int) (domain.Book, error) {
	builder := squirrel.Select("id", "author", "title", "publication_year", "pages", "genre").
		From("books").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return domain.Book{}, fmt.Errorf("db load: failed to build load query: %w", err)
	}

	var b domain.Book
	err = s.db.QueryRow(ctx, sqlStr, args...).Scan(
		&b.ID,
		&b.Author,
		&b.Title,
		&b.PublicationYear,
		&b.Pages,
		&b.Genre,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Book{}, domain.ErrBookNotFound
		}
		return domain.Book{}, fmt.Errorf("db load: query row error: %w", err)
	}

	return b, nil
}

func (s *BookStorage) DeleteWithLog(ctx context.Context, id int) error {

	// begin transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("db delete: failed to create transaction: %w", err)
	}

	// rollback transaction
	defer tx.Rollback(ctx)

	// delete book
	deleteBook := squirrel.Delete("books").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)
	deleteBookSql, args, err := deleteBook.ToSql()
	if err != nil {
		return fmt.Errorf("db delete: failed to build delete query: %w", err)
	}

	result, err := tx.Exec(ctx, deleteBookSql, args...)
	if err != nil {
		return fmt.Errorf("db delete: exec error (delete book): %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrBookNotFound
	}

	// create log
	logID := 0
	createLog := squirrel.Insert("books_log").
		Columns("book_id", "action").
		Values(id, "delete").
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar)

	createLogSql, args, err := createLog.ToSql()
	if err != nil {
		return fmt.Errorf("db delete: failed to build create log query:%w", err)
	}
	row := tx.QueryRow(ctx, createLogSql, args...)
	err = row.Scan(&logID)
	if err != nil {
		return fmt.Errorf("db delete: query row error (create log): %w", err)
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("db delete: commit transaction error: %w", err)
	}
	return nil
}

func (s *BookStorage) UpdateWithLog(ctx context.Context, id int, book domain.Book) error {

	// begin transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("db update: begin transaction error: %w", err)
	}

	// rollback transaction
	defer tx.Rollback(ctx)

	// update book
	updateBook := squirrel.Update("books").
		Set("author", book.Author).
		Set("title", book.Title).
		Set("publication_year", book.PublicationYear).
		Set("pages", book.Pages).
		Set("genre", book.Genre).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	updateBookSql, args, err := updateBook.ToSql()
	if err != nil {
		return fmt.Errorf("db update: failed to build update query: %w", err)
	}

	result, err := tx.Exec(ctx, updateBookSql, args...)
	if err != nil {
		return fmt.Errorf("db update: exec error (update book): %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrBookNotFound
	}

	// create log
	logID := 0
	createLog := squirrel.Insert("books_log").
		Columns("book_id", "action").
		Values(id, "update").
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar)

	createLogSql, args, err := createLog.ToSql()
	if err != nil {
		return fmt.Errorf("db update: failed to build create log query: %w", err)
	}

	row := tx.QueryRow(ctx, createLogSql, args...)
	err = row.Scan(&logID)
	if err != nil {
		return fmt.Errorf("db update: query row error (create log): %w", err)
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("db update: commit transaction error: %w", err)
	}

	return nil
}
