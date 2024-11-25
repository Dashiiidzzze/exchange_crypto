package app

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PairResponse struct {
	Pair_id     int `json:"pair_id"`
	Sale_lot_id int `json:"sale_lot_id"`
	Buy_lot_id  int `json:"buy_lot_id"`
}

// Получение информации о парах:
func HandlePair(w http.ResponseWriter, r *http.Request) {
	var reqBD string = "SELECT * FROM pair"

	response, err := RquestDataBase(reqBD)
	if err != nil {
		return
	}

	fmt.Println(response)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
