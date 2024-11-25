package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
type DeleteOrder struct {
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

		// Проверка наличия ключа пользователя в БД
		var reqUserID string = "SELECT user.user_id FROM user WHERE user.key = '" + userKey + "'"
		userID, err := RquestDataBase(reqUserID)
		if err != nil || userID == "" {
			http.Error(w, "User unauthorized", http.StatusUnauthorized)
			return
		}

		// проверка наличия пары в бд
		var reqPairID string = "SELECT pair.pair_id FROM pair WHERE pair.pair_id = '" + strconv.Itoa(req.PairID) + "'"
		pairID, err1 := RquestDataBase(reqPairID)
		if err1 != nil || pairID == "" {
			http.Error(w, "Pair not found", http.StatusNotFound)
			return
		}

		userID = userID[:len(userID)-2]
		var reqBD string = "INSERT INTO order VALUES ('" + userID + "', '" + strconv.Itoa(req.PairID) + "', '" + strconv.FormatFloat(req.Quantity, 'f', -1, 64) + "', '" + strconv.FormatFloat(req.Price, 'f', -1, 64) + "', '" + req.Type + "', 'open')"
		fmt.Println(reqBD)
		response, err2 := RquestDataBase(reqBD)
		if err2 != nil {
			return
		}
		fmt.Println(response)

		// надо вернуть order_id
		reqBD = "SELECT order.order_id FROM order WHERE order.user_id = '" + userID + "' AND order.closed = 'open'"

		orderIDall, err3 := RquestDataBase(reqBD)
		if err3 != nil {
			return
		}

		orderID := strings.Split(orderIDall, " \n")
		resOrderID, _ := strconv.Atoi(orderID[len(orderID)-2])

		// Формируем и отправляем JSON-ответ клиенту
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CreateOrderResponse{
			OrderID: resOrderID,
		})

	} else if r.Method == "GET" { // получение списка ордеров
		reqBD := "SELECT * FROM order"

		// Имитируем вызов базы данных
		response, err := RquestDataBase(reqBD)
		if err != nil {
			http.Error(w, "Ошибка запроса к базе данных", http.StatusInternalServerError)
			return
		}

		// Преобразуем ответ базы данных в строки
		rows := strings.Split(strings.TrimSpace(response), "\n") // Разделяем строки

		// Массив для хранения ордеров
		var orders []GetOrderResponse

		// Парсим каждую строку
		for _, row := range rows {
			fields := strings.Split(row, " ")
			if len(fields) < 7 {
				continue // Пропускаем строки с недостаточным количеством полей
			}

			// Преобразуем каждое поле и заполняем структуру
			orderID, _ := strconv.Atoi(strings.TrimSpace(fields[0]))
			userID, _ := strconv.Atoi(strings.TrimSpace(fields[1]))
			lotID, _ := strconv.Atoi(strings.TrimSpace(fields[2]))
			quantity, _ := strconv.ParseFloat(strings.TrimSpace(fields[3]), 64)
			orderType := strings.TrimSpace(fields[4])
			price, _ := strconv.ParseFloat(strings.TrimSpace(fields[5]), 64)
			closed := strings.TrimSpace(fields[6])

			order := GetOrderResponse{
				OrderID:  orderID,
				UserID:   userID,
				LotID:    lotID,
				Quantity: quantity,
				Type:     orderType,
				Price:    price,
				Closed:   closed,
			}

			orders = append(orders, order) // Добавляем ордер в массив
		}

		// Устанавливаем заголовки ответа
		w.Header().Set("Content-Type", "application/json")

		// Кодируем массив ордеров в JSON и отправляем клиенту
		json.NewEncoder(w).Encode(orders)
	} else if r.Method == "DELETE" { // удаление ордера
		userKey := r.Header.Get("X-USER-KEY")
		if userKey == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Проверка наличия ключа пользователя в БД
		var reqUserID string = "SELECT user.key FROM user WHERE user.key = '" + userKey + "'"
		userID, err := RquestDataBase(reqUserID)
		if err != nil || userID == "" {
			http.Error(w, "User unauthorized", http.StatusUnauthorized)
			return
		}

		// Парсинг запроса
		var req DeleteOrder
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// проверка является ли пользователь создателем запроса
		userID = userID[:len(userID)-2]
		var reqUserOrder string = "SELECT * FROM order WHERE order.order_id = '" + strconv.Itoa(req.OrderID) + "' AND order.user_id = '" + userID + "'"
		check, err := RquestDataBase(reqUserOrder)
		if err != nil || check == "" {
			http.Error(w, "access error", http.StatusUnauthorized)
			return
		}

		var reqBD string = "DELETE FROM order WHERE order.order_id = '" + strconv.Itoa(req.OrderID) + "'"

		_, err2 := RquestDataBase(reqBD)
		if err2 != nil {
			return
		}

		// Формируем и отправляем JSON-ответ клиенту
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DeleteOrder{
			OrderID: req.OrderID,
		})
	}
}
