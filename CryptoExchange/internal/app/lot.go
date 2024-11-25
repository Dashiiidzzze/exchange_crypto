package app

import (
	"encoding/json"
	"fmt"
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
	if err == nil {
		return
	}
	strResponse := string(response)

	fmt.Println(strResponse)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(strResponse)
}
