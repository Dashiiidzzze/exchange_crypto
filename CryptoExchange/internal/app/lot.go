package app

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type LotResponse struct {
	Lot_id int    `json:"lot_id"`
	Name   string `json:"name"`
}

// Получение информации о лотах
func HandleGetLot(w http.ResponseWriter, r *http.Request) {
	var reqBD string = "SELECT * FROM lot"

	response, err := RquestDataBase(reqBD)
	if err != nil {
		return
	}
	// Преобразуем ответ базы данных в строки
	rows := strings.Split(strings.TrimSpace(response), "\n") // Разделяем строки

	// Массив для хранения ордеров
	var orders []LotResponse

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

		orders = append(orders, order) // Добавляем ордер в массив
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
