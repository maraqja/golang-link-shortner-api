package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func Chain(middlewares ...Middleware) Middleware {
	// Возвращает новую функцию-обработчик, которая принимает и возвращает http.Handler
	return func(next http.Handler) http.Handler {

		// Цикл идет с конца слайса middleware в обратном порядке
		// Это нужно для правильного порядка выполнения middleware:
		// последний добавленный middleware выполнится первым
		for i := len(middlewares) - 1; i >= 0; i-- {

			// Каждый middleware оборачивает предыдущий handler
			// Создается "луковичная" структура обработчиков
			// Например: auth(logging(cors(router)))
			next = middlewares[i](next)
		}

		// Возвращает финальный обработчик, содержащий всю цепочку middleware
		return next
	}

}

/*
// v2:
func Chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        handler = middlewares[i](handler)
    }
    return handler
}


// тогда использовать:
server := http.Server{
    Addr: fmt.Sprintf(":%d", port),
    Handler: middleware.Chain(router,
        middleware.CORS,
        middleware.Logging,
        middleware.Auth,
        middleware.RateLimit,
    ),
}
*/
