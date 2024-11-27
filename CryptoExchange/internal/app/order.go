package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
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
	Type     string  `json:"type"`
	Price    float64 `json:"price"`
	Closed   string  `json:"closed"`
}

// Структура запроса на удаление ордера
type DeleteOrder struct {
	OrderID int `json:"order_id"`
}

// Списание средств
func payByOrder(userID string, pairID int, payMoney float64, price float64, orderType string, spisanie bool) error {
	// Получить информацию о валютной паре
	reqPair := "SELECT * FROM pair WHERE pair.pair_id = '" + strconv.Itoa(pairID) + "'"
	pairData, err := RquestDataBase(reqPair)
	if err != nil || pairData == "" {
		return errors.New("валютная пара не найдена")
	}

	// Разбираем данные пары
	pairFields := strings.Split(pairData, " ") // "1 RUB USD"
	if len(pairFields) < 3 {
		return errors.New("некорректные данные пары")
	}
	firstLotID := pairFields[1]
	secondLotID := pairFields[2]

	// Определяем, с какого счета списывать/начислять
	var lotID string
	if orderType == "buy" {
		lotID = secondLotID // Покупка, списываем вторую валюту
	} else {
		lotID = firstLotID // Продажа, списываем первую валюту
	}

	// Получить баланс пользователя
	reqBalance := "SELECT * FROM user_lot WHERE user_lot.user_id = '" + userID + "' AND user_lot.lot_id = '" + lotID + "'"
	balanceData, err := RquestDataBase(reqBalance)
	if err != nil || balanceData == "" {
		return errors.New("недостаточно средств или запись не найдена")
	}

	balanceFields := strings.Split(balanceData, " ") // "22,2,1,1000"
	if len(balanceFields) < 4 {
		return errors.New("некорректные данные баланса")
	}
	currentBalance, _ := strconv.ParseFloat(balanceFields[3], 64)

	// Проверка средств
	if spisanie && currentBalance < payMoney*price { // добавить домножение
		return errors.New("недостаточно средств")
	}

	// Обновить баланс
	var newBalance float64
	if spisanie {
		newBalance = currentBalance - payMoney*price
	} else {
		newBalance = currentBalance + payMoney*price
	}

	// Удалить старую запись
	reqDelete := "DELETE FROM user_lot WHERE user_lot.user_id = '" + userID + "' AND user_lot.lot_id = '" + lotID + "'"
	_, err = RquestDataBase(reqDelete)
	if err != nil {
		return errors.New("ошибка при обновлении баланса")
	}

	// Вставить новую запись
	reqInsert := fmt.Sprintf("INSERT INTO user_lot VALUES ('%s', '%s', '%.2f')", userID, lotID, newBalance)
	_, err = RquestDataBase(reqInsert)
	if err != nil {
		return errors.New("ошибка при обновлении баланса")
	}

	return nil
}

// поиск уже существующих ордеров для транзакции
func searchOrder(searchUserID string, orderPairID int, orderType string, quantity float64, price float64, types string) (float64, error) {
	var searchOrderTypes string
	if types == "buy" {
		searchOrderTypes = "sell"
	} else {
		searchOrderTypes = "buy"
	}
	// Получить все открытые ордера по паре
	reqOrders := "SELECT * FROM order WHERE order.closed = 'open' AND order.pair_id = '" + strconv.Itoa(orderPairID) + "'"
	fmt.Println("запрос ордера", reqOrders)
	orderData, err := RquestDataBase(reqOrders)
	if err != nil {
		return -1, errors.New("ошибка при поиске ордеров")
	}

	// Разбираем строки
	rows := strings.Split(strings.TrimSpace(orderData), "\n")
	var orders []GetOrderResponse

	for _, row := range rows {
		fields := strings.Split(row, " ")
		if len(fields) < 7 {
			continue
		}

		// Парсим данные
		orderID, _ := strconv.Atoi(fields[0])
		userID, _ := strconv.Atoi(fields[1])
		pairID, _ := strconv.Atoi(fields[2])
		orderQuantity, _ := strconv.ParseFloat(fields[3], 64)
		orderPrice, _ := strconv.ParseFloat(fields[4], 64)
		orderType := fields[5]
		closed := fields[6]

		if (orderType == "buy" && price <= orderPrice || orderType == "sell" && price >= orderPrice) && orderType == searchOrderTypes && strconv.Itoa(userID) != searchUserID {
			orders = append(orders, GetOrderResponse{
				OrderID:  orderID,
				UserID:   userID,
				PairID:   pairID,
				Quantity: orderQuantity,
				Type:     orderType,
				Price:    orderPrice,
				Closed:   closed,
			})
		}
	}

	// Если не найдено подходящих ордеров
	if len(orders) == 0 {
		return quantity, nil
	}

	// Сортировка: для покупки выбираем минимальную цену, для продажи максимальную
	sort.Slice(orders, func(i, j int) bool {
		if orderType == "buy" {
			return orders[i].Price < orders[j].Price
		}
		return orders[i].Price > orders[j].Price
	})

	//собираем подходящие orderID
	totalQuantity := 0.0
	for _, order := range orders {
		if totalQuantity >= quantity { // Если запрос полностью покрыт
			break
		}

		if order.Quantity+totalQuantity > quantity { // Если текущий ордер больше, чем нужно, создаем остаток
			var remainingQuantity float64 = order.Quantity + totalQuantity - quantity
			if remainingQuantity > 0 { // Создаем новый ордер с остатком
				// удаляем ордер
				var forcloseOrderQuery string = "DELETE FROM order WHERE order.order_id = '" + strconv.Itoa(order.OrderID) + "' AND order.closed = 'open'"
				_, err := RquestDataBase(forcloseOrderQuery)
				if err != nil {
					return -1, errors.New("ошибка при закрытии ордера")
				}
				// зачисляем часть денег владельцу ордера (тип ордера противоположный так как произошло завершение транзакции)
				_ = payByOrder(strconv.Itoa(order.UserID), order.PairID, remainingQuantity, order.Price, orderType, false)

				var createOrderQueryClose string = "INSERT INTO order VALUES ('" + strconv.Itoa(order.UserID) + "', '" + strconv.Itoa(order.PairID) + "', '" + strconv.FormatFloat(order.Quantity-remainingQuantity, 'f', -1, 64) + "', '" + strconv.FormatFloat(order.Price, 'f', -1, 64) + "', '" + order.Type + "', 'close')"
				_, err = RquestDataBase(createOrderQueryClose)
				if err != nil {
					return -1, errors.New("ошибка при создании остаточного ордера")
				}

				var createOrderQuery string = "INSERT INTO order VALUES ('" + strconv.Itoa(order.UserID) + "', '" + strconv.Itoa(order.PairID) + "', '" + strconv.FormatFloat(remainingQuantity, 'f', -1, 64) + "', '" + strconv.FormatFloat(order.Price, 'f', -1, 64) + "', '" + order.Type + "', '" + order.Closed + "')"
				_, err = RquestDataBase(createOrderQuery)
				if err != nil {
					return -1, errors.New("ошибка при создании остаточного ордера")
				}
			}

			// Уменьшаем количество в текущем ордере до закрытия
			order.Quantity = quantity - totalQuantity
		} else {
			// удаляем ордер
			var forcloseOrderQuery string = "DELETE FROM order WHERE order.order_id = '" + strconv.Itoa(order.OrderID) + "' AND order.closed = 'open'"
			_, err := RquestDataBase(forcloseOrderQuery)
			if err != nil {
				return -1, errors.New("ошибка при закрытии ордера")
			}

			// зачисляем все деньги владельцу ордера (тип ордера противоположный так как произошло завершение транзакции)
			_ = payByOrder(strconv.Itoa(order.UserID), order.PairID, order.Quantity, order.Price, orderType, false)

			// закрываем ордер
			if order.Quantity+totalQuantity <= quantity {
				var closeOrderQuery string = "INSERT INTO order VALUES ('" + strconv.Itoa(order.UserID) + "', '" + strconv.Itoa(order.PairID) + "', '" + strconv.FormatFloat(order.Quantity, 'f', -1, 64) + "', '" + strconv.FormatFloat(order.Price, 'f', -1, 64) + "', '" + order.Type + "', 'close')"
				_, err := RquestDataBase(closeOrderQuery)
				if err != nil {
					return -1, errors.New("ошибка при закрытии ордера")
				}
			}
		}

		totalQuantity += order.Quantity
	}

	// зачисляем все деньги владельцу исходного ордера (тип ордера противоположный так как произошло завершение транзакции)
	_ = payByOrder(searchUserID, orderPairID, quantity, price, searchOrderTypes, false)

	if totalQuantity < quantity {
		return quantity - totalQuantity, nil
	}

	return 0, nil
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

		// списать средства со счета пользователя
		payErr := payByOrder(userID, req.PairID, req.Quantity, req.Price, req.Type, true)
		if payErr != nil {
			http.Error(w, "Not enough funds", http.StatusPaymentRequired)
			return
		}

		// здесь вставить поиск подходящего ордера на рокупку продажу, если нашелся, начисляем новые средства
		newQuant, searchError := searchOrder(userID, req.PairID, req.Type, req.Quantity, req.Price, req.Type)
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
			_, err := RquestDataBase(closeOrderQuery)
			if err != nil {
				return
			}
			status = "open"
		} else {
			status = "open"
		}
		var reqBD string = "INSERT INTO order VALUES ('" + userID + "', '" + strconv.Itoa(req.PairID) + "', '" + strconv.FormatFloat(newQuant, 'f', -1, 64) + "', '" + strconv.FormatFloat(req.Price, 'f', -1, 64) + "', '" + req.Type + "', '" + status + "')"
		_, err2 := RquestDataBase(reqBD)
		if err2 != nil {
			return
		}
		// надо вернуть order_id (предполагается, что это последний ордер, добавленный в бд)
		reqBD = "SELECT order.order_id FROM order WHERE order.user_id = '" + userID + "' AND order.closed = '" + status + "'"

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
			pairID, _ := strconv.Atoi(strings.TrimSpace(fields[2]))
			quantity, _ := strconv.ParseFloat(strings.TrimSpace(fields[3]), 64)
			orderType := strings.TrimSpace(fields[4])
			price, _ := strconv.ParseFloat(strings.TrimSpace(fields[5]), 64)
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
		var reqUserOrder string = "SELECT * FROM order WHERE order.order_id = '" + strconv.Itoa(req.OrderID) + "' AND order.user_id = '" + userID + "' AND order.closed = 'open'"
		check, err2 := RquestDataBase(reqUserOrder)
		if err2 != nil || check == "" {
			http.Error(w, "access error", http.StatusUnauthorized)
			return
		}
		balanceFields := strings.Split(check, " ")
		if len(balanceFields) < 7 {
			return
		}

		var reqBD string = "DELETE FROM order WHERE order.order_id = '" + strconv.Itoa(req.OrderID) + "'"

		_, err3 := RquestDataBase(reqBD)
		if err3 != nil {
			return
		}

		// вернуть деньги обратно на счет пользователю
		floatQuant, _ := strconv.ParseFloat(strings.TrimSpace(balanceFields[3]), 64)
		floatPrice, _ := strconv.ParseFloat(strings.TrimSpace(balanceFields[4]), 64)
		num, _ := strconv.Atoi(balanceFields[2])
		payErr := payByOrder(userID, num, floatQuant, floatPrice, balanceFields[5], false)
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
