package main

import (
	"CryptoExchange/internal/app"
	"CryptoExchange/internal/config"
	"log"
	"net/http"
)

func main() {
	// формирование таблицы с парами
	pairList := config.ConfigRead()
	app.Init(pairList)

	// Регистрируем обработчики
	http.HandleFunc("/user", app.HandleCreateUser)    // POST
	http.HandleFunc("/order", app.HandleOrder)        // POST, GET, DELETE
	http.HandleFunc("/lot", app.HandleGetLot)         // GET,
	http.HandleFunc("/pair", app.HandlePair)          // GET
	http.HandleFunc("/balance", app.HandleGetBalance) // GET

	// Запускаем сервер на порту 8080
	log.Println("Сервер запущен на порту 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
