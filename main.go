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
	http.HandleFunc("/hello", hello)
	port:=8081
	fmt.Printf("Server is listening on port %d ...", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil) // начинаем обработку TCP-соединений на порту 8080 и передаем в него дефолтный обработчик (ServeMux - принимает запросы и передает их в соответствующие обработчики)
}