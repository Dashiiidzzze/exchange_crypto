package app

import (
	"CryptoExchange/internal/logic"
	"CryptoExchange/internal/requestDB"
	"encoding/json"
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
	PairID   int     `json:"lot_id"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
	Type     string  `json:"type"`
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

		// Проверка наличия ключа пользователя в БД
		var reqUserID string = "SELECT user.user_id FROM user WHERE user.key = '" + userKey + "'"
		userID, err := requestDB.RquestDataBase(reqUserID)
		if err != nil || userID == "" {
			http.Error(w, "User unauthorized", http.StatusUnauthorized)
			return
		}
		userID = userID[:len(userID)-2]

		// проверка наличия пары в бд
		var reqPairID string = "SELECT pair.pair_id FROM pair WHERE pair.pair_id = '" + strconv.Itoa(req.PairID) + "'"
		pairID, err1 := requestDB.RquestDataBase(reqPairID)
		if err1 != nil || pairID == "" {
			http.Error(w, "Pair not found", http.StatusNotFound)
			return
		}

		// списать средства со счета пользователя
		payErr := logic.PayByOrder(userID, req.PairID, req.Quantity, req.Price, req.Type, true)
		if payErr != nil {
			http.Error(w, "Not enough funds", http.StatusPaymentRequired)
			return
		}

		// здесь вставить поиск подходящего ордера на рокупку продажу, если нашелся, начисляем новые средства
		newQuant, searchError := logic.SearchOrder(userID, req.PairID, req.Type, req.Quantity, req.Price, req.Type)
		if searchError != nil {
			http.Error(w, "Not enough orders", http.StatusNotFound)
			return
		}
		// создаем ордер
		status := ""
		if newQuant == 0 {
			status = "close"
			newQuant = req.Quantity
		} else if newQuant != req.Quantity {
			// вносим в базу уже закрытый ордер (точнее его часть)
			var closeOrderQuery string = "INSERT INTO order VALUES ('" + userID + "', '" + strconv.Itoa(req.PairID) + "', '" + strconv.FormatFloat(req.Quantity, 'f', -1, 64) + "', '" + strconv.FormatFloat(req.Price, 'f', -1, 64) + "', '" + req.Type + "', 'close')"
			_, err := requestDB.RquestDataBase(closeOrderQuery)
			if err != nil {
				return
			}
			status = "open"
		} else {
			status = "open"
		}
		var reqBD string = "INSERT INTO order VALUES ('" + userID + "', '" + strconv.Itoa(req.PairID) + "', '" + strconv.FormatFloat(newQuant, 'f', -1, 64) + "', '" + strconv.FormatFloat(req.Price, 'f', -1, 64) + "', '" + req.Type + "', '" + status + "')"
		_, err2 := requestDB.RquestDataBase(reqBD)
		if err2 != nil {
			return
		}
		// надо вернуть order_id (предполагается, что это последний ордер, добавленный в бд)
		reqBD = "SELECT order.order_id FROM order WHERE order.user_id = '" + userID + "' AND order.closed = '" + status + "'"

		orderIDall, err3 := requestDB.RquestDataBase(reqBD)
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
		reqBD := "SELECT * FROM order WHERE order.closed = 'open'"

		// Имитируем вызов базы данных
		response, err := requestDB.RquestDataBase(reqBD)
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
			pairID, _ := strconv.Atoi(strings.TrimSpace(fields[2]))
			quantity, _ := strconv.ParseFloat(strings.TrimSpace(fields[3]), 64)
			orderType := strings.TrimSpace(fields[5])
			price, _ := strconv.ParseFloat(strings.TrimSpace(fields[4]), 64)
			closed := strings.TrimSpace(fields[6])

			order := GetOrderResponse{
				OrderID:  orderID,
				UserID:   userID,
				PairID:   pairID,
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
		var reqUserID string = "SELECT user.user_id FROM user WHERE user.key = '" + userKey + "'"
		userID, err := requestDB.RquestDataBase(reqUserID)
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
		var reqUserOrder string = "SELECT * FROM order WHERE order.order_id = '" + strconv.Itoa(req.OrderID) + "' AND order.user_id = '" + userID + "' AND order.closed = 'open'"
		check, err2 := requestDB.RquestDataBase(reqUserOrder)
		if err2 != nil || check == "" {
			http.Error(w, "access error", http.StatusUnauthorized)
			return
		}
		balanceFields := strings.Split(check, " ")
		if len(balanceFields) < 7 {
			return
		}

		var reqBD string = "DELETE FROM order WHERE order.order_id = '" + strconv.Itoa(req.OrderID) + "'"

		_, err3 := requestDB.RquestDataBase(reqBD)
		if err3 != nil {
			return
		}

		// вернуть деньги обратно на счет пользователю
		floatQuant, _ := strconv.ParseFloat(strings.TrimSpace(balanceFields[3]), 64)
		floatPrice, _ := strconv.ParseFloat(strings.TrimSpace(balanceFields[4]), 64)
		num, _ := strconv.Atoi(balanceFields[2])
		payErr := logic.PayByOrder(userID, num, floatQuant, floatPrice, balanceFields[5], false)
		if payErr != nil {
			http.Error(w, "Not enough funds", http.StatusPaymentRequired)
			return
		}

		// Формируем и отправляем JSON-ответ клиенту
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DeleteOrder{
			OrderID: req.OrderID,
		})
	}
}
