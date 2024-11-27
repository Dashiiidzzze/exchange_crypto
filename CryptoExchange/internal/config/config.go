package config

import (
	"encoding/json"
	"log"
	"os"
)

// Структура для хранения данных из JSON
type ConfigStruct struct {
	Lots         []string `json:"lots"`
	DatabaseIP   string   `json:"database_ip"`
	APIPort      int      `json:"api_port"`
	DatabasePort int      `json:"database_port"`
}

func ConfigRead() ([]string, string, int, int) {
	file, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Не удалось открыть файл: %v", err)
	}
	//fmt.Println(string(file))

	var config ConfigStruct
	if err := json.Unmarshal(file, &config); err != nil {
		log.Fatalf("Ошибка при парсинге JSON: %v", err)
	}

	return config.Lots, config.DatabaseIP, config.APIPort, config.DatabasePort
}
