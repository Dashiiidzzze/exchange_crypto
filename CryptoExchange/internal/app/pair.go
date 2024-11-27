package app

import (
	"CryptoExchange/internal/requestDB"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type PairResponse struct {
	Pair_id     int `json:"pair_id"`
	Sale_lot_id int `json:"sale_lot_id"`
	Buy_lot_id  int `json:"buy_lot_id"`
}

// Получение информации о парах:
func HandlePair(w http.ResponseWriter, r *http.Request) {
	var reqBD string = "SELECT * FROM pair"

	response, err := requestDB.RquestDataBase(reqBD)
	if err != nil {
		return
	}
	// Преобразуем ответ базы данных в строки
	rows := strings.Split(strings.TrimSpace(response), "\n") // Разделяем строки

	// Массив для хранения
	var pairs []LotResponse

	// Парсим каждую строку
	for _, row := range rows {
		fields := strings.Split(row, " ")
		if len(fields) < 2 {
			continue // Пропускаем строки с недостаточным количеством полей
		}

		// Преобразуем каждое поле и заполняем структуру
		lotID, _ := strconv.Atoi(strings.TrimSpace(fields[0]))
		name := strings.TrimSpace(fields[1])

		order := LotResponse{
			Lot_id: lotID,
			Name:   name,
		}

		pairs = append(pairs, order) // Добавляем ордер в массив
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pairs)
}
