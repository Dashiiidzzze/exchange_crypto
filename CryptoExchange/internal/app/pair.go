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
	var pairs []PairResponse

	// Парсим каждую строку
	for _, row := range rows {
		fields := strings.Split(row, " ")
		if len(fields) < 3 {
			continue // Пропускаем строки с недостаточным количеством полей
		}

		// Преобразуем каждое поле и заполняем структуру
		pairID, _ := strconv.Atoi(strings.TrimSpace(fields[0]))
		sale_lot_id, _ := strconv.Atoi(strings.TrimSpace(fields[1]))
		buy_lot_id, _ := strconv.Atoi(strings.TrimSpace(fields[2]))

		pa := PairResponse{
			Pair_id:     pairID,
			Sale_lot_id: sale_lot_id,
			Buy_lot_id:  buy_lot_id,
		}

		pairs = append(pairs, pa) // Добавляем ордер в массив
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pairs)
}
