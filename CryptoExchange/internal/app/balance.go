package app

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	var reqUserID string = "SELECT user.key FROM user WHERE user.key = '" + userKey + "'"
	userID := RquestDataBase(reqUserID)
	if userID == nil {
		http.Error(w, "User unauthorized", http.StatusUnauthorized)
		return
	}

	var reqBD string = "SELECT user_lot.lot_id user_lot.quantity FROM user_lot WHERE user_lot.user_id = '" + string(userID) + "'"

	response := RquestDataBase(reqBD)
	if response == nil {
		return
	}
	strResponse := string(response)

	fmt.Println(strResponse)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(strResponse)
}
