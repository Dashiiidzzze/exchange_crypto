package app

import (
	"encoding/json"
	"net/http"
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
