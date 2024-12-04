package main

import (
	"context"
	"fmt"
	"link-shortner-api/configs"
	"link-shortner-api/internal/auth"
	"link-shortner-api/internal/link"
	"link-shortner-api/internal/user"
	"link-shortner-api/pkg/db"
	"link-shortner-api/pkg/middleware"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()                                       // создаем базовый контекст
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 2*time.Second) // из него создаем контекст с таймаутом
	defer cancel()                                                    // если выполнили быстрее таймаута, то отменим контекст

	done := make(chan struct{})

	go func() {
		time.Sleep(3 * time.Second)
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("Done task")
	case <-ctxWithTimeout.Done():
		fmt.Println("Timeout")
	}
}

func main2() {
	port := 8081
	config := configs.LoadConfig()

	db := db.NewDb(config)
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)

	authService := auth.NewAuthService(userRepository)

	router := http.NewServeMux() // создаем новый ServeMux

	// Handler
	auth.NewAuthHandler(router, &auth.AuthHandlerDependencies{
		Config:      config,
		AuthService: authService,
	})
	link.NewLinkHandler(router, &link.LinkHandlerDependencies{
		LinkRepository: linkRepository,
	})

	// Middlewares
	middlewareStack := middleware.Chain( // вызываются в таком же порядке
		middleware.CORS,
		middleware.Logging,
	)

	server := http.Server{ // конфигурируем сервер через структуру
		Addr:    fmt.Sprintf(":%d", port),
		Handler: middlewareStack(router),
	}

	fmt.Printf("Server is listening on port %d ...", port)
	server.ListenAndServe()
}
