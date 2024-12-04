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
)

func main() {
	type key string
	const EmailKey key = "email"                                         // ключ для хранения email в контексте
	ctx := context.Background()                                          // создаем контекст
	ctxWithValue := context.WithValue(ctx, EmailKey, "test@example.com") // добавляем значение в контекст

	userEmail, ok := ctxWithValue.Value(EmailKey).(string) // получаем значение из контекста и преобразовываем в строку
	// В Go оператор приведения типа .(T) используется для приведения значения к определенному типу T
	if ok {
		fmt.Println(userEmail)
	} else {
		fmt.Println("No Value")
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
