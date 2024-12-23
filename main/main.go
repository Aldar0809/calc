package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Aldar0809/calc/api"
)

func main() {
	// Регистрируем обработчик
	http.HandleFunc("/api/v1/calculate", api.CalculateHandler)

	// Запускаем сервер
	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
