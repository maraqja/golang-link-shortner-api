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

func tickOperation(ctx context.Context) {
	ticker := time.NewTicker(200 * time.Millisecond) // содержит канал, который можно слушать (каждый указанный период будет слать время тика)
	for {
		select {
		case tick := <-ticker.C:
			fmt.Println(tick)
		case <-ctx.Done():
			fmt.Println("cancel")
			return
		}

	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go tickOperation(ctx)

	time.Sleep(2 * time.Second)
	cancel()
	time.Sleep(1 * time.Second) // нужно для того чтобы дождаться выполнения горутины с отменой

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
