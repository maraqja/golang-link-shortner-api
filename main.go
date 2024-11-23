package main

import (
	"fmt"
	"net/http"
)

// фукнция обработчик роута /hello
// Входные параметры:
// http.ResponseWriter - куда будем писать
// *http.Request - указатель на содержимое запроса
func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Hello")
} 

func main() {

	/*
	ServeMux:
		Принимает входящие HTTP-запросы
		Сопоставляет URL-пути с зарегистрированными обработчиками (handlers)
		Направляет запросы в соответствующие обработчики
	*/


	/*
	DefaultServeMux является глобальной переменной в пакете http, это означает что:
		Любая часть программы может регистрировать в нем обработчики через http.HandleFunc()
		Все зарегистрированные маршруты хранятся в одном месте
		При масштабировании приложения разные пакеты могут случайно перезаписать маршруты друг друга
	*/

	// Поэтому нужно использовать собственный serveMux для каждого из модулей
	port:=8081
	
	router := http.NewServeMux() // создаем новый ServeMux
	router.HandleFunc("/hello", hello)
	server := http.Server{ // конфигурируем сервер через структуру
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}


	fmt.Printf("Server is listening on port %d ...", port)
	server.ListenAndServe()
}