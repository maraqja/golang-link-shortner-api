package main

import (
	"fmt"
	"link-shortner-api/configs"
	"link-shortner-api/internal/auth"
	"link-shortner-api/internal/link"
	"link-shortner-api/pkg/db"
	"link-shortner-api/pkg/middleware"
	"net/http"
)

func main() {
	port := 8081
	config := configs.LoadConfig()

	db := db.NewDb(config)
	linkRepository := link.NewLinkRepository(db)

	router := http.NewServeMux() // создаем новый ServeMux

	// Handler
	auth.NewAuthHandler(router, &auth.AuthHandlerDependencies{
		Config: config,
	})
	link.NewLinkHandler(router, &link.LinkHandlerDependencies{
		LinkRepository: linkRepository,
	})

	server := http.Server{ // конфигурируем сервер через структуру
		Addr:    fmt.Sprintf(":%d", port),
		Handler: middleware.Logging(router),
	}

	fmt.Printf("Server is listening on port %d ...", port)
	server.ListenAndServe()
}
