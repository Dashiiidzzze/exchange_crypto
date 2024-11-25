package app

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type BalanceResponse struct {
	Lot_id   int     `json:"lot_id"`
	Quantity float64 `json:"quantity"`
}

func HandleGetBalance(w http.ResponseWriter, r *http.Request) {
	userKey := r.Header.Get("X-USER-KEY")
	// Проверка наличия заголовка X-USER-KEY, проверить есть ли такой пользователь
	if userKey == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var reqUserID string = "SELECT user.user_id user.key FROM user WHERE user.key = '" + userKey + "'"
	userIDandKey, err := RquestDataBase(reqUserID)
	if err != nil {
		http.Error(w, "User unauthorized", http.StatusUnauthorized)
		return
	}
	userID := strings.Split(userIDandKey, " ")

	var reqBD string = "SELECT user_lot.lot_id user_lot.quantity FROM user_lot WHERE user_lot.user_id = '" + userID[0] + "'"

	response, err2 := RquestDataBase(reqBD)
	if err2 != nil {
		return
	}
	// Преобразуем ответ базы данных в строки
	rows := strings.Split(strings.TrimSpace(response), "\n") // Разделяем строки

	// Массив для хранения ордеров
	var orders []BalanceResponse

	// Парсим каждую строку
	for _, row := range rows {
		fields := strings.Split(row, " ")
		if len(fields) < 2 {
			continue // Пропускаем строки с недостаточным количеством полей
		}

		// Преобразуем каждое поле и заполняем структуру
		lotID, _ := strconv.Atoi(strings.TrimSpace(fields[0]))
		quantity, _ := strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)

		order := BalanceResponse{
			Lot_id:   lotID,
			Quantity: quantity,
		}

		orders = append(orders, order) // Добавляем ордер в массив
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
