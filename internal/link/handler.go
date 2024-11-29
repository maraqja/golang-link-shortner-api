package link

import (
	"fmt"
	"link-shortner-api/pkg/request"
	"link-shortner-api/pkg/response"
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
		body, err := request.HandleBody[LinkCreateRequest](&w, req)
		if err != nil {
			return // тк респонс с ошибкой уже отправлен в HandleBody, возвращаемся
		}
		// по сути код отсюда и до отправки респонса должен быть в отдельном сервисе
		link := NewLink(body.Url)
		for { // в бесконечном цикле проверяем, что ссылки с таким хешем нет
			existedLink, _ := handler.LinkRepository.GetByHash(link.Hash)
			if existedLink == nil { // если уникальный хэш сгенерился, то выходим из цикла
				break
			}
			link.GenerateHash() // если неуникальный хэш сгенерился - генерируем новый
		}

		createdLink, err := handler.LinkRepository.Create(link)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		response.ReturnJSON(w, http.StatusCreated, createdLink)
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
		hash := req.PathValue("hash")
		link, err := handler.LinkRepository.GetByHash(hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, req, link.Url, http.StatusTemporaryRedirect) // если все ок, то редиректим на ссылку
	}
}
