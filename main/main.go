package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"calc"
)

// RequestBody представляет тело запроса
type RequestBody struct {
	Expression string `json:"expression"`
}

// ResponseBody представляет тело ответа
type ResponseBody struct {
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

// CalculateHandler обрабатывает запросы на вычисление выражений
func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что запрос является POST
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Декодируем тело запроса
	var reqBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Некорректное тело запроса", http.StatusBadRequest)
		return
	}

	// Вычисляем выражение
	result, err := calc.Calc(reqBody.Expression)
	if err != nil {
		// Обрабатываем ошибки
		var statusCode int
		var errorMessage string
		switch err.Error() {
		case "пустое выражение", "некорректное выражение", "недопустимый токен", "непарные скобки":
			statusCode = http.StatusUnprocessableEntity
			errorMessage = "Expression is not valid"
		default:
			statusCode = http.StatusInternalServerError
			errorMessage = "Internal server error"
		}
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(ResponseBody{Error: errorMessage})
		return
	}

	// Возвращаем результат
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseBody{Result: result})
}

func main() {
	// Регистрируем обработчик
	http.HandleFunc("/api/v1/calculate", CalculateHandler)

	// Запускаем сервер
	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
