package restapi

import (
	apiv1 "task4/internal/delivery/rest_api/apiV1"
	"task4/internal/usecase"

	"github.com/gin-gonic/gin"
)

type router struct {
	router      *gin.Engine
	bookHandler *apiv1.BookHandler
}

func NewRouter(r *gin.Engine, bu usecase.BookUseCase) *router {
	return &router{
		router:      r,
		bookHandler: apiv1.NewBookHandler(&bu),
	}
}

func (r *router) Register() {
	r.router.GET("/book", r.bookHandler.List)
	r.router.GET("/book/:id", r.bookHandler.Load)
	r.router.POST("/book", r.bookHandler.Create)
	r.router.DELETE("/book/:id", r.bookHandler.Delete)
	r.router.PATCH("/book/:id", r.bookHandler.Update)
}
