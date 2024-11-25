package main

import (
	"CryptoExchange/internal/app"
	"CryptoExchange/internal/config"
	"log"
	"net/http"
	"strconv"
)

func main() {
	// формирование таблицы с парами
	pairList, _, port, _ := config.ConfigRead()
	app.Init(pairList)

	// Регистрируем обработчики
	http.HandleFunc("/user", app.HandleCreateUser)    // POST
	http.HandleFunc("/order", app.HandleOrder)        // POST, GET, DELETE
	http.HandleFunc("/lot", app.HandleGetLot)         // GET,
	http.HandleFunc("/pair", app.HandlePair)          // GET
	http.HandleFunc("/balance", app.HandleGetBalance) // GET

	// Запускаем сервер на порту 8080
	log.Println("Сервер запущен на порту " + strconv.Itoa(port) + " ...")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))

}
