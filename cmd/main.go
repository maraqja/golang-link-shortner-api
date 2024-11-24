package main

import (
	"fmt"
	"link-shortner-api/configs"
	"link-shortner-api/internal/hello"
	"net/http"
)



func main() {

	conf := configs.LoadConfig()

	port := 8081

	router := http.NewServeMux() // создаем новый ServeMux
	hello.NewHelloHandler(router)
	server := http.Server{ // конфигурируем сервер через структуру
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}


	fmt.Printf("Server is listening on port %d ...", port)
	server.ListenAndServe()
}