package responce

import (
	"task4/internal/domain"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func SendError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, Response{
		Success: false,
		Message: message,
	})
}

func SendSuccess(c *gin.Context, code int, data any) {
	c.JSON(code, Response{
		Success: true,
		Data:    data,
	})
}

type BookResponce struct {
	ID              int    `json:"id"`
	Author          string `json:"author"`
	Title           string `json:"title"`
	PublicationYear int    `json:"publication_year"`
	Pages           int    `json:"pages"`
	Genre           string `json:"genre"`
}

func (r *BookResponce) FromDomain(b domain.Book) {
	r.ID = b.ID
	r.Author = b.Author
	r.Title = b.Title
	r.PublicationYear = b.PublicationYear
	r.Pages = b.Pages
	r.Genre = b.Genre
}

type BookListResponce struct {
	Books []BookResponce `json:"books"`
	Total int            `json:"total"`
}
