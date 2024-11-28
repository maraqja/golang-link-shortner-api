package link

import (
	"fmt"
	"net/http"
)

type LinkHandlerDependencies struct {
	*LinkRepository
}

type LinkHandler struct {
	*LinkRepository
}

func NewLinkHandler(router *http.ServeMux, dependencies *LinkHandlerDependencies) {
	handler := &LinkHandler{
		LinkRepository: dependencies.LinkRepository,
	}
	router.Handle("POST /link", handler.Create())
	router.Handle("PATCH /link/{id}", handler.Update())
	router.Handle("DELETE /link/{id}", handler.Delete())
	router.Handle("GET /{hash}", handler.GoTo())
}

func (handler *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		handler.LinkRepository.Create(&Link{})
	}
}

func (handler *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id := req.PathValue("id")
		fmt.Println(id)
	}
}

func (handler *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

	}
}

func (handler *LinkHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

	}
}
