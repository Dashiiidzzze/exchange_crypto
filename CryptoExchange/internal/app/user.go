package app

import (
	"CryptoExchange/internal/config"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
)

// Структура для запроса создания пользователя
type CreateUserRequest struct {
	Username string `json:"username"`
}

// Структура ответа при создании пользователя
type CreateUserResponse struct {
	Key string `json:"key"`
}

// Генерация уникального ключа для пользователя
func generateUserKey(username string) string {
	n, _ := rand.Int(rand.Reader, big.NewInt(100000)) // Генерация случайного числа
	hash := md5.Sum([]byte(username + strconv.Itoa(int(n.Int64()))))
	return hex.EncodeToString(hash[:])
}

// Генерация уникального ключа для пользователя
func assetGen(userKey string) {
	var reqBDcheck string = "SELECT user.user_id FROM user WHERE user.key = '" + userKey + "'"
	response, err := RquestDataBase(reqBDcheck)
	if err != nil {
		return
	}
	response = response[:len(response)-2]
	lots := config.ConfigRead()
	for i := 0; i < len(lots); i++ {
		var reqBDsearch string = "SELECT lot.lot_id FROM lot WHERE lot.name = '" + lots[i] + "'"
		lotID, err2 := RquestDataBase(reqBDsearch)
		if err2 != nil {
			return
		}
		lotID = lotID[:len(lotID)-2]
		var reqBD string = "INSERT INTO user_lot VALUES ('" + response + "', '" + lotID + "', '1000')"
		_, err3 := RquestDataBase(reqBD)
		if err3 != nil {
			return
		}
	}
}

// Функция для создания пользователя
func HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	// Парсинг JSON-запроса от клиента
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}
	fmt.Println(req)

	userKey := generateUserKey(req.Username)

	var reqBD string = "INSERT INTO user VALUES ('" + req.Username + "', '" + userKey + "')"

	_, err := RquestDataBase(reqBD)
	if err != nil {
		return
	}
	// генерация активов пользователя
	assetGen(userKey)

	// Формируем и отправляем JSON-ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CreateUserResponse{Key: string(userKey)})
}
