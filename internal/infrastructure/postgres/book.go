package postgres

import (
	"context"
	"errors"
	"fmt"
	"task4/internal/domain"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
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

func (s *BookStorage) Create(ctx context.Context, data domain.Book) (id int, err error) {
	builder := squirrel.Insert("books").
		Columns("author", "title", "publication_year", "pages", "genre").
		Values(data.Author, data.Title, data.PublicationYear, data.Pages, data.Genre).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar)

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return id, errors.New("failed to generate SQL statement: " + err.Error())
	}

	row := s.db.QueryRow(ctx, sqlStr, args...)
	err = row.Scan(&id)
	if err != nil {
		return id, fmt.Errorf("db insert book error: %w", err)
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
		return domain.Book{}, fmt.Errorf("failed to build load query: %w", err)
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
		if errors.Is(err, pgx.ErrNoRows) || err.Error() == "sql: no rows in result set" {
			return domain.Book{}, domain.ErrBookNotFound
		}
		return domain.Book{}, fmt.Errorf("db load book error: %w", err)
	}

	return b, nil
}

func (s *BookStorage) Delete(ctx context.Context, id int) error {
	builder := squirrel.Delete("books").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)
	sql, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	result, err := s.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("db delete book error: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrBookNotFound
	}

	return nil
}

func (s *BookStorage) Update(ctx context.Context, id int, book domain.Book) error {
	builder := squirrel.Update("books").
		Set("author", book.Author).
		Set("title", book.Title).
		Set("publication_year", book.PublicationYear).
		Set("pages", book.Pages).
		Set("genre", book.Genre).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	result, err := s.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("db update book error: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrBookNotFound
	}

	return nil
}
