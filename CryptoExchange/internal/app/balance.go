package app

import (
	"encoding/json"
	"net/http"
	"strings"
)

type BalanceResponse struct {
	Lot_id   int    `json:"lot_id"`
	Quantity string `json:"quantity"`
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
