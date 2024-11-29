package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"time"
)

const apiURL = "http://localhost:8080"

type OrderRequest struct {
	PairID   int     `json:"pair_id"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
	Type     string  `json:"type"`
}

type Pair struct {
	PairID      int `json:"pair_id"`
	Sale_lot_id int `json:"sale_lot_id"`
	Buy_lot_id  int `json:"buy_lot_id"`
}

type UserResponse struct {
	Key string `json:"key"`
}

// создание HTTP запроса post (для создания пользователя и ордеров)
func postRequest(endpoint string, payload any, apiKey string) ([]byte, error) {
	data, _ := json.Marshal(payload) // кодирование в json
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

// создание HTTP запроса get (для получения пар)
func getRequest(endpoint string) ([]byte, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", apiURL+endpoint, nil)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func main() {
	rand.NewSource(time.Now().UnixNano())

	// Создаем пользователя
	user := struct {
		Username string `json:"username"`
	}{"random_bot"}
	resp, _ := postRequest("/user", user, "")
	var userResponse UserResponse
	json.Unmarshal(resp, &userResponse)
	apiKey := userResponse.Key

	user = struct {
		Username string `json:"username"`
	}{"random_bot_2"}
	resp, _ = postRequest("/user", user, "")
	var userResponse2 UserResponse
	json.Unmarshal(resp, &userResponse2)
	apiKeyTWO := userResponse2.Key

	// Получаем доступные пары
	pairsResp, _ := getRequest("/pair")
	var pairs []Pair
	json.Unmarshal(pairsResp, &pairs)
	fmt.Println("Доступные пары:", pairs)

	// Бесконечный цикл для работы случайного бота
	for {
		pair := pairs[rand.Intn(len(pairs))] // Случайная пара
		order := OrderRequest{
			PairID:   pair.PairID,
			Quantity: math.Round((rand.Float64()*100+1)*100) / 100, // Округление до 2 знаков
			Price:    math.Round((rand.Float64()*10+1)*100) / 100,  // Округление до 2 знаков
			Type:     []string{"buy", "sell"}[rand.Intn(2)],
		}
		_, err := postRequest("/order", order, []string{apiKey, apiKeyTWO}[rand.Intn(2)])
		if err != nil {
			fmt.Println("Ошибка создания ордера:", err)
		} else {
			fmt.Println("Создан случайный ордер:", order)
		}
		time.Sleep(1 * time.Second)
	}
}
