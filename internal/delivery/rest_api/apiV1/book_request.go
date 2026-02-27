package apiv1

import "task4/internal/domain"

type ListBooksRequest struct {
	Author    string `form:"author"`
	Title     string `form:"title"`
	Genre     string `form:"genre"`
	YearFrom  int    `form:"year_from"`
	YearTo    int    `form:"year_to"`
	PagesFrom int    `form:"pages_from"`
	PagesTo   int    `form:"pages_to"`
	Limit     uint64 `form:"limit,default=10"`
	Offset    uint64 `form:"offset,default=0"`
}

func (l *ListBooksRequest) ToDomain() domain.BookFilter {
	return domain.BookFilter{
		Author:    l.Author,
		Title:     l.Title,
		Genre:     l.Genre,
		YearFrom:  l.YearFrom,
		YearTo:    l.YearTo,
		PagesFrom: l.PagesFrom,
		PagesTo:   l.PagesTo,
		Limit:     l.Limit,
		Offset:    l.Offset,
	}
}

type CreateBookRequest struct {
	Author          string `json:"author"`
	Title           string `json:"title"`
	PublicationYear int    `json:"publication_year"`
	Pages           int    `json:"pages"`
	Genre           string `json:"genre"`
}

func (r *CreateBookRequest) ToDomain() domain.Book {
	return domain.Book{
		Author:          r.Author,
		Title:           r.Title,
		PublicationYear: r.PublicationYear,
		Pages:           r.Pages,
		Genre:           r.Genre,
	}
}

type UpdateBookRequest struct {
	Author          string `json:"author"`
	Title           string `json:"title"`
	PublicationYear int    `json:"publication_year"`
	Pages           int    `json:"pages"`
	Genre           string `json:"genre"`
}

func (r *UpdateBookRequest) ToDomain() domain.Book {
	return domain.Book{
		Author:          r.Author,
		Title:           r.Title,
		PublicationYear: r.PublicationYear,
		Pages:           r.Pages,
		Genre:           r.Genre,
	}
}
