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

// // списание средств со счета пользователя
// func payByOrder(userID string, pairID int, payMoney float64, spisanie bool) error {
// 	var reqBDlots string = "SELECT pair.firat_lot_id pair.second_lot_id FROM pair WHERE pair.pair_id = '" + strconv.Itoa(pairID) + "'"

// 	lots, err := RquestDataBase(reqBDlots)
// 	if err != nil {
// 		return errors.New("не удалось подключиться к базе данных")
// 	}
// 	allLots := strings.Split(lots, " ")

// 	// проверка есть ли у пользователя деньги
// 	var reqBDselect string = "SELECT user_lot.quantity FROM user_lot WHERE user_lot.user_id = '" + userID + "' AND user_lot.lot_id = '" + allLots[1] + "'"

// 	money, err2 := RquestDataBase(reqBDselect)
// 	if err2 != nil {
// 		return errors.New("не удалось подключиться к базе данных")
// 	}
// 	money = money[:len(money)-2]
// 	intMoney, _ := strconv.ParseFloat(strings.TrimSpace(money), 64)

// 	if spisanie {
// 		if intMoney < payMoney {
// 			return errors.New("недостаточно средств")
// 		}
// 	}

// 	// списание (обновление актива на меньшую сумму)
// 	var reqBDdel string = "DELETE FROM user_lot WHERE user_lot.user_id = '" + userID + "'"

// 	_, err3 := RquestDataBase(reqBDdel)
// 	if err3 != nil {
// 		return errors.New("не удалось подключиться к базе данных")
// 	}

// 	var reqBD string
// 	if spisanie {
// 		reqBD = "INSERT INTO user_lot VALUES ('" + userID + "', '" + allLots[1] + "', '" + strconv.FormatFloat(intMoney-payMoney, 'f', -1, 64) + "')"
// 	} else {
// 		reqBD = "INSERT INTO user_lot VALUES ('" + userID + "', '" + allLots[1] + "', '" + strconv.FormatFloat(intMoney+payMoney, 'f', -1, 64) + "')"
// 	}
// 	_, err3 = RquestDataBase(reqBD)
// 	if err3 != nil {
// 		return errors.New("не удалось подключиться к базе данных")
// 	}

// 	return nil
// }

// func searchOrder(order_id string, user_id string, pair_id string, quantyty string, price string, types string) {
// 	var reqBD string
// 	if types == "buy" {
// 		reqBD = "SELECT * FROM order WHERE order.pair_id = '" + pair_id + "' AND order.type = 'sell' AND order.closed = 'open'"
// 	} else {
// 		reqBD = "SELECT * FROM order WHERE order.pair_id = '" + pair_id + "' AND order.type = 'buy' AND order.closed = 'open'"
// 	}
// 	response, err := RquestDataBase(reqBD)
// 	if err != nil {
// 		return
// 	}

// 	rows := strings.Split(response, "\n") // Разделяем строки
// 	// Парсим каждую строку
// 	for _, row := range rows {
// 		fields := strings.Split(row, " ")
// 		if len(fields) < 7 {
// 			continue // Пропускаем строки с недостаточным количеством полей
// 		}
// 		var orderIDSearch string
// 		if types == "buy" && price >= rows[4] && quantyty == rows[3] {
// 			orderIDSearch = rows[0]
// 		} else if types == "sell" && price <= rows[4] && quantyty == rows[3] {
// 			orderIDSearch = rows[0]
// 		}
// 		if orderIDSearch!= "" {
// 			var check string = "DELETE FROM order WHERE order.order_id = '" + orderIDSearch + "'"

// 			_, err2 := RquestDataBase(check)
// 			if err2 != nil {
// 				return
// 			}

// 			var check2 string = "DELETE FROM order WHERE order.order_id = '" + order_id + "'"

// 			_, err2 = RquestDataBase(check2)
// 			if err2 != nil {
// 				return
// 			}

// 			var check3 string = "INSERT INTO order VALUES ('" + orderIDSearch + "', '" + rows[1] + "', '" + rows[2] + "', '" + rows[3] + "', '" + rows[4] + "', '" + rows[5] + "', 'close'"

// 			_, err2 = RquestDataBase(check3)
// 			if err2 != nil {
// 				return
// 			}

// 			var check4 string = "INSERT INTO order VALUES ('" + order_id + "', '" + user_id + "', '" + strconv.Itoa(pair_id) + "', '" + strconv.FormatFloat(quantity, 'f', -1, 64) + "', '" + strconv.FormatFloat(price, 'f', -1, 64) + "', 'close')"

// 			_, err2 = RquestDataBase(check4)
// 			if err2 != nil {
// 				return
// 			}
// 			payByOrder(rows[1], rows[2], rows[3], spisanie bool, aftersearch bool)
// 		}
// 	}

// 	// Массив для хранения ордеров
// 	var orders []GetOrderResponse
// }

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
		userID = userID[:len(userID)-2]

		// проверка наличия пары в бд
		var reqPairID string = "SELECT pair.pair_id FROM pair WHERE pair.pair_id = '" + strconv.Itoa(req.PairID) + "'"
		pairID, err1 := RquestDataBase(reqPairID)
		if err1 != nil || pairID == "" {
			http.Error(w, "Pair not found", http.StatusNotFound)
			return
		}

		// здесь списать средства со счета пользователя
		// payErr := payByOrder(userID, req.PairID, req.Quantity*req.Price, true)
		// if payErr != nil {
		// 	http.Error(w, "Not enough funds", http.StatusPaymentRequired)
		// 	return
		// }

		// здесь вставить поиск подходящего ордера на рокупку продажу, если нашелся, начисляем новые средства
		// searchError := searchOrder()
		// if searchError != nil {
		// 	http.Error(w, "Not enough funds", http.StatusPaymentRequired)
		// 	return
		// }

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
		// order := strings.Split(check, " ")

		var reqBD string = "DELETE FROM order WHERE order.order_id = '" + strconv.Itoa(req.OrderID) + "'"

		_, err2 := RquestDataBase(reqBD)
		if err2 != nil {
			return
		}

		// здесь вернуть деньги обратно на счет пользователю
		// здесь списать средства со счета пользователя
		// floatQuant, _ := strconv.ParseFloat(strings.TrimSpace(order[3]), 64)
		// floatPrice, _ := strconv.ParseFloat(strings.TrimSpace(order[4]), 64)
		// num, _ := strconv.Atoi(order[2])
		// payErr := payByOrder(userID, num, floatQuant*floatPrice, true)
		// if payErr != nil {
		// 	http.Error(w, "Not enough funds", http.StatusPaymentRequired)
		// 	return
		// }

		// Формируем и отправляем JSON-ответ клиенту
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DeleteOrder{
			OrderID: req.OrderID,
		})
	}
}
