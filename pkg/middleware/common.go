package middleware

import "net/http"

type WrapperWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *WrapperWriter) WriteHeader(statusCode int) { // переопределяем метод, чтобы сохранить статус код
	w.ResponseWriter.WriteHeader(statusCode)
	w.StatusCode = statusCode
}
