package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type CreateOrderRequest struct {
	PairID   int     `json:"pair_id"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
	Type     string  `json:"type"`
}

// Структура ответа при создании ордера
type CreateOrderResponse struct {
	OrderID int `json:"order_id"`
}

type GetOrderResponse struct {
	OrderID  int     `json:"order_id"`
	UserID   int     `json:"user_id"`
	LotID    int     `json:"lot_id"`
	Quantity float64 `json:"quantity"`
	Type     string  `json:"type"`
	Price    float64 `json:"price"`
	Closed   string  `json:"closed"`
}

// Структура запроса на удаление ордера
type DeleteOrderRequest struct {
	OrderID int `json:"order_id"`
}

// работа с ордерами
func HandleOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" { // создание ордера
		userKey := r.Header.Get("X-USER-KEY")
		// Проверка наличия заголовка X-USER-KEY, проверить есть ли такой пользователь
		if userKey == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Парсинг JSON-запроса
		var req CreateOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		fmt.Println(req)

		var reqUserID string = "SELECT user.key FROM user WHERE user.key = '" + userKey + "'"
		userID := RquestDataBase(reqUserID)
		if userID == nil {
			http.Error(w, "User unauthorized", http.StatusUnauthorized)
			return
		}

		stringUserID := string(userID)
		stringUserID = stringUserID[:len(stringUserID)-4]
		//var reqBD string = "INSERT INTO order VALUES ('" + string(req.PairID) + "', '" + strconv.FormatFloat(req.Quantity, 'f', -1, 64) + "', '" + req.Price + "', '" + req.Type + "')"
		var reqBD string = "INSERT INTO order VALUES ('" + stringUserID + "', '" + strconv.Itoa(req.PairID) + "', '" + strconv.FormatFloat(req.Quantity, 'f', -1, 64) + "', '" + strconv.FormatFloat(req.Price, 'f', -1, 64) + "', '" + req.Type + "', 'open')"
		fmt.Println(reqBD)
		response := RquestDataBase(reqBD)
		if response == nil {
			return
		}
		strResponse := string(response)
		fmt.Println(strResponse)
		// надо вернуть order_id!!!!!!
		// reqBD  = "SELECT order.order_id FROM order WHERE order."

		// response = RquestDataBase(reqBD)
		// if response == nil {
		// 	return
		// }

		// Формируем и отправляем JSON-ответ клиенту
		w.Header().Set("Content-Type", "application/json")
		//json.NewEncoder(w).Encode(CreateOrderResponse({Key: string(userKey)}))
		json.NewEncoder(w).Encode(strResponse)

	} else if r.Method == "GET" { // получение списка ордеров
		var reqBD string = "SELECT * FROM order"

		response := RquestDataBase(reqBD)
		if response == nil {
			return
		}
		strResponse := string(response)
		fmt.Println(strResponse)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(strResponse)

	} else if r.Method == "DELETE" { // удаление ордера
		userKey := r.Header.Get("X-USER-KEY")
		if userKey == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Парсинг запроса
		var req DeleteOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		var reqBD string = "DELETE FROM order WHERE order.order_id = '" + strconv.Itoa(req.OrderID) + "'"

		response := RquestDataBase(reqBD)
		if response == nil {
			return
		}
		strResponse := string(response)
		fmt.Println(strResponse)
	}
}
