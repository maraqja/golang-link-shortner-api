package response

import (
	"encoding/json"
	"net/http"
)

func ReturnJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Add("Content-Type", "application/json") // добавляем хедер
	w.WriteHeader(statusCode) // устанавливаем статус-код
	json.NewEncoder(w).Encode(data) // создаем json-encoder (преобразует структуры в json формат), который будет писать напрямую в респонс

}