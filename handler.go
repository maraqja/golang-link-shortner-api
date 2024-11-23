package main

import (
	"fmt"
	"net/http"
)

// HelloHandler структура для группировки обработчиков
type HelloHandler struct{}

// Hello возвращает http.HandlerFunc, которая автоматически реализует интерфейс http.Handler
// http.HandlerFunc - это тип, который реализует ServeHTTP под капотом
func (h *HelloHandler) Hello() http.HandlerFunc {
    // Возвращаем функцию с сигнатурой ServeHTTP(ResponseWriter, *Request)
    return func(w http.ResponseWriter, req *http.Request) {
        fmt.Println("Hello")
    }
}

// NewHelloHandler регистрирует обработчик в маршрутизаторе
// router.Handle ожидает http.Handler интерфейс, который требует метод ServeHTTP
// http.HandlerFunc автоматически преобразует нашу функцию в http.Handler
func NewHelloHandler(router *http.ServeMux) {
    handler := &HelloHandler{}
    router.Handle("/hello", handler.Hello()) // Hello() возвращает HandlerFunc, которая реализует ServeHTTP
}