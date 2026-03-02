package apiv1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"task4/internal/domain"
	mockUseCase "task4/internal/usecase/mock"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookCreate(t *testing.T) {
	mockUC := new(mockUseCase.BookUseCaseMock)
	testHandler := &BookHandler{
		bu: mockUC,
	}

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/book", testHandler.Create)

	t.Run("success", func(t *testing.T) {
		expectedID := 123

		expectedBook := domain.Book{
			Author:          "test_author",
			Title:           "test_title",
			PublicationYear: 1999,
			Pages:           123,
			Genre:           "test_genre",
		}

		reqBook := CreateBookRequest{
			Author:          "test_author",
			Title:           "test_title",
			PublicationYear: 1999,
			Pages:           123,
			Genre:           "test_genre",
		}

		mockUC.On("Create", mock.Anything, expectedBook).Return(expectedID, nil).Once()

		body, _ := json.Marshal(reqBook)
		req, _ := http.NewRequest(http.MethodPost, "/book", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("invalid_body", func(t *testing.T) {

		invalidBook := []byte(`{"title": "broken json"`)

		body, _ := json.Marshal(invalidBook)
		req, _ := http.NewRequest(http.MethodPost, "/book", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal_error", func(t *testing.T) {

		expectedBook := domain.Book{
			Author:          "test_author",
			Title:           "test_title",
			PublicationYear: 1999,
			Pages:           123,
			Genre:           "test_genre",
		}

		reqBook := CreateBookRequest{
			Author:          "test_author",
			Title:           "test_title",
			PublicationYear: 1999,
			Pages:           123,
			Genre:           "test_genre",
		}

		mockUC.On("Create", mock.Anything, expectedBook).Return(0, fmt.Errorf("internal_error")).Once()

		body, _ := json.Marshal(reqBook)
		req, _ := http.NewRequest(http.MethodPost, "/book", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUC.AssertExpectations(t)
	})
}

func TestBookLoad(t *testing.T) {
	mockUC := new(mockUseCase.BookUseCaseMock)
	testHandler := &BookHandler{
		bu: mockUC,
	}

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/book/:id", testHandler.Load)

	t.Run("success", func(t *testing.T) {
		expectedID := 123

		expectedBook := domain.Book{
			ID:              expectedID,
			Author:          "test_author",
			Title:           "test_title",
			PublicationYear: 1999,
			Pages:           123,
			Genre:           "test_genre",
		}

		mockUC.On("Load", mock.Anything, expectedID).Return(expectedBook, nil).Once()

		url := fmt.Sprintf("/book/%d", expectedID)

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("incorrect_id", func(t *testing.T) {

		badId := "very_very_bad_id"
		badUrl := fmt.Sprintf("/book/%s", badId)

		req, _ := http.NewRequest(http.MethodGet, badUrl, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not_found_id", func(t *testing.T) {

		expectedID := 123

		expectedBook := domain.Book{}

		mockUC.On("Load", mock.Anything, expectedID).Return(expectedBook, domain.ErrBookNotFound).Once()

		url := fmt.Sprintf("/book/%d", expectedID)

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("internal_err", func(t *testing.T) {

		expectedID := 123

		expectedBook := domain.Book{}

		mockUC.On("Load", mock.Anything, expectedID).Return(expectedBook, fmt.Errorf("definitely internal server error")).Once()

		url := fmt.Sprintf("/book/%d", expectedID)

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUC.AssertExpectations(t)
	})
}

func TestBookList(t *testing.T) {
	mockUC := new(mockUseCase.BookUseCaseMock)
	testHandler := &BookHandler{
		bu: mockUC,
	}

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/book", testHandler.List)

	expectedBooks := []domain.Book{
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

	t.Run("success_without_filters", func(t *testing.T) {

		emptyFilter := domain.BookFilter{
			Limit: 10, //default value for limit
		}

		mockUC.On("List", mock.Anything, emptyFilter).Return(expectedBooks, len(expectedBooks), nil).Once()

		url := "/book"

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("success_with_filters", func(t *testing.T) {

		filter := domain.BookFilter{
			Limit:     10, //default value for limit
			PagesFrom: 120,
			YearFrom:  1900,
		}

		mockUC.On("List", mock.Anything, filter).Return(expectedBooks, len(expectedBooks), nil).Once()

		url := "/book?pages_from=120&year_from=1900"

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("bad_query", func(t *testing.T) {

		url := "/book?pages_from=one_hundred&year_from=NAN"

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not_found", func(t *testing.T) {

		mockUC.On("List", mock.Anything, mock.Anything).Return(nil, 0, nil).Once()

		url := "/book"

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("internal_error", func(t *testing.T) {

		mockUC.On("List", mock.Anything, mock.Anything).Return(nil, 0, fmt.Errorf("definitely internal err")).Once()

		url := "/book"

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUC.AssertExpectations(t)
	})
}

func TestBookDelete(t *testing.T) {
	mockUC := new(mockUseCase.BookUseCaseMock)
	testHandler := &BookHandler{
		bu: mockUC,
	}

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.DELETE("/book/:id", testHandler.Delete)

	t.Run("success", func(t *testing.T) {
		expectedID := 123

		mockUC.On("Delete", mock.Anything, expectedID).Return(nil).Once()

		url := fmt.Sprintf("/book/%d", expectedID)

		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("invalid_id", func(t *testing.T) {

		badId := "badid"
		url := fmt.Sprintf("/book/%s", badId)

		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("not_found", func(t *testing.T) {
		expectedID := 123

		mockUC.On("Delete", mock.Anything, expectedID).Return(domain.ErrBookNotFound).Once()

		url := fmt.Sprintf("/book/%d", expectedID)

		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("internal_err", func(t *testing.T) {
		expectedID := 123

		mockUC.On("Delete", mock.Anything, expectedID).Return(fmt.Errorf("definitely internal err")).Once()

		url := fmt.Sprintf("/book/%d", expectedID)

		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUC.AssertExpectations(t)
	})
}

func TestBookUpdate(t *testing.T) {
	mockUC := new(mockUseCase.BookUseCaseMock)
	testHandler := &BookHandler{
		bu: mockUC,
	}

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.PATCH("/book/:id", testHandler.Update)

	t.Run("success", func(t *testing.T) {
		expectedID := 123
		expectedBook := domain.Book{
			Author:          "test_author",
			Title:           "test_title",
			PublicationYear: 1999,
			Pages:           123,
			Genre:           "test_genre",
		}
		// в метод айдишник идет из запроса, а не из данных объекта
		// подразумевается, что айдишник объекта в бд изменить нельзя
		mockUC.On("Update", mock.Anything, expectedID, expectedBook).Return(nil).Once()

		url := fmt.Sprintf("/book/%d", expectedID)

		reqBook := UpdateBookRequest{
			Author:          "test_author",
			Title:           "test_title",
			PublicationYear: 1999,
			Pages:           123,
			Genre:           "test_genre",
		}

		body, _ := json.Marshal(reqBook)
		req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("invalid_id", func(t *testing.T) {

		badId := "badid"
		url := fmt.Sprintf("/book/%s", badId)

		reqBook := UpdateBookRequest{
			Author:          "test_author",
			Title:           "test_title",
			PublicationYear: 1999,
			Pages:           123,
			Genre:           "test_genre",
		}

		body, _ := json.Marshal(reqBook)

		req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("invalid_body", func(t *testing.T) {
		expectedID := 123

		url := fmt.Sprintf("/book/%d", expectedID)

		badBody := []byte(`{
			"data": im very bad body,
		}`)

		req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(badBody))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("not_found", func(t *testing.T) {
		expectedID := 123

		mockUC.On("Update", mock.Anything, expectedID, mock.Anything).Return(domain.ErrBookNotFound).Once()

		url := fmt.Sprintf("/book/%d", expectedID)

		reqBook := UpdateBookRequest{
			Author:          "test_author",
			Title:           "test_title",
			PublicationYear: 1999,
			Pages:           123,
			Genre:           "test_genre",
		}

		body, _ := json.Marshal(reqBook)
		req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("internal_err", func(t *testing.T) {
		expectedID := 123

		mockUC.On("Update", mock.Anything, expectedID, mock.Anything).Return(fmt.Errorf("definitely internal err")).Once()

		url := fmt.Sprintf("/book/%d", expectedID)

		reqBook := UpdateBookRequest{
			Author:          "test_author",
			Title:           "test_title",
			PublicationYear: 1999,
			Pages:           123,
			Genre:           "test_genre",
		}

		body, _ := json.Marshal(reqBook)
		req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUC.AssertExpectations(t)
	})
}
