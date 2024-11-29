package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

const apiURL = "http://localhost:8080"

// Структуры данных
type Pair struct {
	PairID    int `json:"pair_id"`
	SaleLotID int `json:"sale_lot_id"`
	BuyLotID  int `json:"buy_lot_id"`
}

type Lot struct {
	LotID int    `json:"lot_id"`
	Name  string `json:"name"`
}

type Order struct {
	OrderID  int     `json:"order_id"`
	UserID   int     `json:"user_id"`
	PairID   int     `json:"lot_id"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
	Type     string  `json:"type"`
	Closed   string  `json:"closed"`
}

type Balance struct {
	LotID    int     `json:"lot_id"`
	Quantity float64 `json:"quantity"`
}

type UserResponse struct {
	Key string `json:"key"`
}

type OrderRequest struct {
	PairID   int     `json:"pair_id"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
	Type     string  `json:"type"`
}

// Функция для POST-запросов
func postRequest(endpoint string, payload any, apiKey string) ([]byte, error) {
	data, _ := json.Marshal(payload)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", apiURL+endpoint, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("X-USER-KEY", apiKey)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// Функция для GET-запросов
func getRequest(endpoint string, apiKey string) ([]byte, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", apiURL+endpoint, nil)
	if apiKey != "" {
		req.Header.Set("X-USER-KEY", apiKey)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func main() {
	// Создаем пользователя
	user := struct {
		Username string `json:"username"`
	}{"user"}
	resp, _ := postRequest("/user", user, "")
	var userResponse UserResponse
	json.Unmarshal(resp, &userResponse)
	apiKey := userResponse.Key

	// Получаем доступные пары
	pairsResp, _ := getRequest("/pair", apiKey)
	var pairs []Pair
	json.Unmarshal(pairsResp, &pairs)

	// Определяем ID лота RUB
	lotsResp, _ := getRequest("/lot", apiKey)
	var lots []Lot
	json.Unmarshal(lotsResp, &lots)
	var rubLotID int
	for _, lot := range lots {
		if lot.Name == "RUB" {
			rubLotID = lot.LotID
			break
		}
	}

	// Фильтруем пары, где RUB участвует
	var rubPairs []Pair
	for _, pair := range pairs {
		if pair.SaleLotID == rubLotID || pair.BuyLotID == rubLotID {
			rubPairs = append(rubPairs, pair)
		}
	}

	// Бесконечный цикл работы
	for {
		// Получаем текущий баланс
		balanceResp, _ := getRequest("/balance", apiKey)
		var balances []Balance
		json.Unmarshal(balanceResp, &balances)
		var rubBalance float64
		for _, balance := range balances {
			if balance.LotID == rubLotID {
				rubBalance = balance.Quantity
				break
			}
		}
		fmt.Println("баланс:", rubBalance)

		// Получаем список ордеров
		ordersResp, _ := getRequest("/order", apiKey)
		var orders []Order
		json.Unmarshal(ordersResp, &orders)

		var minSell float64 = 100000000
		var pairIDSell int
		var quantitySell float64
		var maxBuy float64 = -100000000
		var pairIDBuy int
		var quantityBuy float64
		var averagePrice float64
		var check int = 0
		for _, order := range orders {
			for _, pair := range rubPairs {
				if order.PairID == pair.PairID {
					if order.Type == "sell" {
						if order.Price < minSell {
							minSell = order.Price
							pairIDSell = order.PairID
							quantitySell = order.Quantity
						}
					} else if order.Type == "buy" {
						if order.Price > maxBuy {
							maxBuy = order.Price
							pairIDBuy = order.PairID
							quantityBuy = order.Quantity
						}
					}
					averagePrice += order.Price
					check++
					break
				}
			}
		}
		if check != 0 {
			averagePrice = averagePrice / float64(check)
			// выгодно выгодно купить так как цена маленькая
			var order OrderRequest
			if math.Abs(averagePrice-minSell) > math.Abs(averagePrice-maxBuy) && minSell != 100000000 {
				order = OrderRequest{
					PairID:   pairIDSell,
					Quantity: quantitySell,
					Price:    minSell,
					Type:     "buy",
				}
				// выгодно продать так как цена большая
			} else {
				order = OrderRequest{
					PairID:   pairIDBuy,
					Quantity: quantityBuy,
					Price:    maxBuy,
					Type:     "sell",
				}
			}
			_, err := postRequest("/order", order, apiKey)
			if err == nil {
				fmt.Printf("выставлен лот: %v\n", order)
				rubBalance -= minSell
			}
		} else {
			fmt.Println("нет ордеров на продажу или покупку")
		}
		time.Sleep(5 * time.Second) // Пауза между итерациями
	}
}
