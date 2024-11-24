package main

import (
	"fmt"
	"link-shortner-api/configs"
	"link-shortner-api/internal/auth"
	"net/http"
)



func main() {

	conf := configs.LoadConfig()

	port := 8081

	router := http.NewServeMux() // создаем новый ServeMux
	auth.NewAuthHandler(router, &auth.AuthHandlerDependencies{
		Config: conf,
	}) 
	server := http.Server{ // конфигурируем сервер через структуру
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}


	fmt.Printf("Server is listening on port %d ...", port)
	server.ListenAndServe()
}