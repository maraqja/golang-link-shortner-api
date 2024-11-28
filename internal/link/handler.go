package link

import (
	"net/http"
)

type LinkHandlerDependencies struct {
}

type LinkHandler struct {
}

func NewLinkHandler(router *http.ServeMux, dependencies *LinkHandlerDependencies) {
	handler := &LinkHandler{}
	router.Handle("POST /link", handler.Create())
	router.Handle("PATCH /link/{id}", handler.Update())
	router.Handle("DELETE /link/{id}", handler.Delete())
	router.Handle("GET /{alias}", handler.GoTo())
}

func (handler *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

	}
}

func (handler *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

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
