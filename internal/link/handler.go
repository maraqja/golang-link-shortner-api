package link

import (
	"errors"
	"fmt"
	"link-shortner-api/configs"
	"link-shortner-api/pkg/middleware"
	"link-shortner-api/pkg/request"
	"link-shortner-api/pkg/response"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type LinkHandlerDependencies struct {
	*LinkRepository
	Config *configs.Config
}

type LinkHandler struct {
	*LinkRepository
}

func NewLinkHandler(router *http.ServeMux, dependencies *LinkHandlerDependencies) {
	handler := &LinkHandler{
		LinkRepository: dependencies.LinkRepository,
	}
	router.HandleFunc("POST /link", handler.Create())
	router.Handle("PATCH /link/{id}", middleware.IsAuthed(handler.Update(), dependencies.Config)) // Здесь уже необходимо использовать Handle, тк middleware.IsAuthed() возвращает http.Handler
	router.HandleFunc("DELETE /link/{id}", handler.Delete())
	router.HandleFunc("GET /{hash}", handler.GoTo())
	router.Handle("GET /link", middleware.IsAuthed(handler.GetAll(), dependencies.Config))
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
			return // обязательно делаем return, иначе будет продолжать выполняться код ниже (и будет несколько строк выведено в респонсе)
		}
		response.ReturnJSON(w, http.StatusCreated, createdLink)
	}
}

func (handler *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		email, ok := req.Context().Value(middleware.ContextEmailKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized) // или 500 ???
			return
		}
		fmt.Println(email)

		body, err := request.HandleBody[LinkUpdateRequest](&w, req)
		if err != nil {
			return
		}
		idString := req.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		link := &Link{ // не создаем через NewLink конструктор тк хотим дать возможность менять сам hash
			Model: gorm.Model{ID: uint(id)},
			Url:   body.Url,
			Hash:  body.Hash,
		}

		updatedLink, err := handler.LinkRepository.Update(link)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response.ReturnJSON(w, http.StatusOK, updatedLink)
	}
}

func (handler *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		idString := req.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// // чекаем есть ли запись (но более эффективно просто проверять кол-во заафекченных при удалении строк)
		// _, err = handler.GetById(uint(id))
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusNotFound)
		// 	return
		// }

		// удаляем запись
		err = handler.LinkRepository.Delete(uint(id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.ReturnJSON(w, http.StatusOK, nil)
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

func (handler *LinkHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		limit, err := strconv.Atoi(req.URL.Query().Get("limit")) // получаем query и конвертируем в int (конвертация может упасть с ошибкой)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid limit: %s", err.Error()), http.StatusBadRequest)
			return
		}
		offset, err := strconv.Atoi(req.URL.Query().Get("offset")) // получаем query и конвертируем в int (конвертация может упасть с ошибкой)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid offset: %s", err.Error()), http.StatusBadRequest)
			return
		}
		links, err := handler.LinkRepository.GetAll(limit, offset)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка при получении ссылок: %s:", err.Error()), http.StatusInternalServerError)
			return
		}
		count, err := handler.LinkRepository.GetCount()
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка при получении общего количества ссылок: %s:", err.Error()), http.StatusInternalServerError)
			return
		}

		response.ReturnJSON(w, http.StatusOK, GetAllLinksResponse{
			Link:  links,
			Count: count,
		})

	}
}
