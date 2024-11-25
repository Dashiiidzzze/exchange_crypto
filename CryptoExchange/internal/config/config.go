package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Структура для хранения данных из JSON
type ConfigStruct struct {
	Lots         []string `json:"lots"`
	DatabaseIP   string   `json:"database_ip"`
	DatabasePort int      `json:"database_port"`
}

func ConfigRead() []string {

	fmt.Println("Config")
	file, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Не удалось открыть файл: %v", err)
	}
	//fmt.Println(string(file))

	var config ConfigStruct
	if err := json.Unmarshal(file, &config); err != nil {
		log.Fatalf("Ошибка при парсинге JSON: %v", err)
	}

	return config.Lots
}
