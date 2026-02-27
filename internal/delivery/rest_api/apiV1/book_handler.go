package apiv1

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"task4/internal/delivery/rest_api/responce"
	"task4/internal/domain"

	"github.com/gin-gonic/gin"
)

type BookUseCase interface {
	Create(ctx context.Context, data domain.Book) (int, error)
	Update(ctx context.Context, id int, data domain.Book) error
	List(ctx context.Context, f domain.BookFilter) ([]domain.Book, int, error)
	Load(ctx context.Context, id int) (domain.Book, error)
	Delete(ctx context.Context, id int) error
}

type BookHandler struct {
	bu BookUseCase
}

func NewBookHandler(bu BookUseCase) *BookHandler {
	return &BookHandler{
		bu: bu,
	}
}

func (h *BookHandler) List(c *gin.Context) {
	req := ListBooksRequest{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		responce.SendError(c, http.StatusBadRequest, "bind query err")
		return
	}
	filter := req.ToDomain()
	books, total, err := h.bu.List(c.Request.Context(), filter)
	if err != nil {
		responce.SendError(c, http.StatusInternalServerError, "internal err")
		return
	}
	if len(books) == 0 {
		responce.SendError(c, http.StatusNotFound, "books not found")
		return
	}

	data := []responce.BookResponce{}
	for _, b := range books {
		temp := responce.BookResponce{}
		temp.FromDomain(b)
		data = append(data, temp)
	}

	resp := responce.BookListResponce{
		Books: data,
		Total: total,
	}

	responce.SendSuccess(c, http.StatusOK, resp)
}

func (h *BookHandler) Load(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responce.SendError(c, http.StatusBadRequest, "invalid id")
		return
	}
	book, err := h.bu.Load(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrBookNotFound) {
			responce.SendError(c, http.StatusNotFound, "book not found")
			return
		}
		responce.SendError(c, http.StatusInternalServerError, "internal err")
		return
	}

	resp := responce.BookResponce{}
	resp.FromDomain(book)

	responce.SendSuccess(c, http.StatusOK, resp)
}

func (h *BookHandler) Create(c *gin.Context) {
	req := CreateBookRequest{}
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		responce.SendError(c, http.StatusBadRequest, "invalid body")
		return
	}
	data := req.ToDomain()
	id, err := h.bu.Create(c.Request.Context(), data)
	if err != nil {
		responce.SendError(c, http.StatusInternalServerError, "internal err")
		return
	}

	responce.SendSuccess(c, http.StatusOK, id)
}

func (h *BookHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responce.SendError(c, http.StatusBadRequest, "invalid id")
		return
	}
	err = h.bu.Delete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrBookNotFound) {
			responce.SendError(c, http.StatusNotFound, "book not found")
			return
		}
		responce.SendError(c, http.StatusInternalServerError, "internal err")
		return
	}

	responce.SendSuccess(c, http.StatusOK, id)
}

func (h *BookHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responce.SendError(c, http.StatusBadRequest, "invalid id")
		return
	}
	req := UpdateBookRequest{}
	err = c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		responce.SendError(c, http.StatusBadRequest, "invalid body")
		return
	}
	book := req.ToDomain()
	err = h.bu.Update(c.Request.Context(), id, book)
	if err != nil {
		if errors.Is(err, domain.ErrBookNotFound) {
			responce.SendError(c, http.StatusNotFound, "book not found")
			return
		}
		responce.SendError(c, http.StatusInternalServerError, "internal err")
		return
	}

	responce.SendSuccess(c, http.StatusOK, id)
}
