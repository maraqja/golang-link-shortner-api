package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		wrapper := &WrapperWriter{ // нужно сделать обертку над ResponseWriter, чтобы сохранить статус код и иметь к нему доступ
			ResponseWriter: w,
			StatusCode:     http.StatusOK, // это просто дефолтное значение, которое будет перезаписано в случае, если статус код будет изменен
		}
		next.ServeHTTP(wrapper, req)                                                 // вызываем следующий хэндлер в цепочке
		log.Println(wrapper.StatusCode, req.Method, req.URL.Path, time.Since(start)) // логируем с помощью встроенной библиотеки

	})
}
