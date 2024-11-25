package app

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

// запрос к базе данных
// func RquestDataBase(req CreateUserRequest) []byte {
func RquestDataBase(req string) (string, error) {
	// Устанавливаем TCP-соединение с базой данных на порту 7432
	//conn, err := net.Dial("tcp", "localhost:7432")
	conn, err := net.Dial("tcp", "db:7432")
	if err != nil {
		//http.Error(w, "Не удалось подключиться к базе данных", http.StatusInternalServerError)
		fmt.Println("Не удалось подключиться к базе данных", http.StatusInternalServerError)
		return "", errors.New("не удалось подключиться к базе данных")
	}
	defer conn.Close() // Закрываем соединение по завершении

	// Отправляем запрос в базу данных
	fmt.Fprintf(conn, req+"\n") // Добавляем перевод строки, если база ожидает его

	// Читаем ответ от базы данных
	response, err := io.ReadAll(conn)
	if err != nil {
		//http.Error(w, "Ошибка при чтении ответа от базы данных", http.StatusInternalServerError)
		fmt.Println("Ошибка при чтении ответа от базы данных", http.StatusInternalServerError)
		return "", errors.New("не удалось подключиться к базе данных")
	}
	str := string(response)
	return str, nil
}
