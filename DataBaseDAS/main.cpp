#include <iostream>
#include <filesystem>
#include <thread>
#include <vector>
#include <string.h>
#include <arpa/inet.h>
#include <unistd.h>

#include "include/mapDas.h"
#include "include/vectorDas.h"
#include "func/header.h"

using namespace std; 

// Парсит и выполняет SQL-запросы
void parsingQuery(const string& query, SchemaInfo& schemaData, int clientSocket) {
    MyVector<string>* words = Split(query, ' ');
    string result;
    if (words->data[0] == "SELECT") {
        try {
            ParsingSelect(*words, schemaData, clientSocket);
        } catch (const exception& err) {
            result = string("Error: ") + err.what() + "\n";
            //cerr << err.what() << endl;
        }
    
    } else if (words->data[0] == "INSERT" && words->data[1] == "INTO") {
        try {
            ParsingInsert(*words, schemaData);
            result = "successful insert\n";
        } catch (const exception& err) {
            result = string("Error: ") + err.what() + "\n";
            //cerr << err.what() << endl;
        }
    
    } else if (words->data[0] == "DELETE" && words->data[1] == "FROM") {
        try {
            ParsingDelete(*words, schemaData);
            result = "successful deletion\n";
        } catch (const exception& err) {
            result = string("Error: ") + err.what() + "\n";
            //cerr << err.what() << endl;
        }
        
    } else { 
        //cout << "Unknown command" << endl;
        result = "Unknown command\n";
    }

    send(clientSocket, result.c_str(), result.size(), 0);
}

// чтение имени файла и пути к нему
//bool inputNames(int clientSocket, string& jsonFileName, SchemaInfo& schemaData) {
bool inputNames(string& jsonFileName, SchemaInfo& schemaData) {
    schemaData.jsonStructure = CreateMap<string, MyVector<string>*>(10, 50);

    // Проверка существования файла
    try {
        if (!filesystem::exists(schemaData.filepath + "/" + jsonFileName)) {
            // string errorMsg = "Error: JSON file not found\n";
            // send(clientSocket, errorMsg.c_str(), errorMsg.size(), 0);   // Отправка сообщения об ошибке
            // close(clientSocket);
            cerr << "Error: JSON file not found" << endl;
            return false;
        } else {
            // Чтение структуры JSON-файла
            ReadJsonFile(jsonFileName, schemaData);
            //string successMsg = "File successfully loaded\n";
            //send(clientSocket, successMsg.c_str(), successMsg.size(), 0);
            return true;
        }
    } catch (const exception& err) {
        // string errorMsg = "Error: " + string(err.what()) + "\n";
        // send(clientSocket, errorMsg.c_str(), errorMsg.size(), 0);
        // close(clientSocket);
        throw;
        return false;
    }
    return false;
}

// Функция для чтения SQL-запросов клиента
//void handleClient(int clientSocket) {
void handleClient(int clientSocket, SchemaInfo& schemaData) {
    // string jsonFileName = "schema.json";
    // SchemaInfo schemaData;
    // //ввод имени файла и пути
    // if (!inputNames(clientSocket, jsonFileName, schemaData)) {
    //     close(clientSocket);
    //     return;
    // }
    // Обработка SQL-запросов клиента
    // while (true) {
    //     char buffer[1024];
    //     //send(clientSocket, "Введите SQL запрос или \"q\" для выхода\n >>> ", 68, 0);
    //     memset(buffer, 0, sizeof(buffer));
    //     ssize_t bytesRead = read(clientSocket, buffer, sizeof(buffer) - 1);
        
    //     if (bytesRead <= 0) {
    //         cerr << "Connection closed by client or error occurred." << endl;
    //         close(clientSocket);
    //         break;
    //     }

    //     string query = string(buffer);
    //     query.erase(query.find_last_not_of("\r\n") + 1); // Удаление символов конца строки

    //     if (query == "q") {
    //         close(clientSocket);
    //         break;
    //     }

    //     // Парсинг запроса
    //     parsingQuery(query, schemaData, clientSocket);
    // }

    char buffer[1024];
    //send(clientSocket, "Введите SQL запрос или \"q\" для выхода\n >>> ", 68, 0);
    memset(buffer, 0, sizeof(buffer));
    ssize_t bytesRead = read(clientSocket, buffer, sizeof(buffer) - 1);
    if (bytesRead <= 0) {
        cerr << "Connection closed by client or error occurred." << endl;
        close(clientSocket);
        return;
    }
    string query = string(buffer);
    query.erase(query.find_last_not_of("\r\n") + 1); // Удаление символов конца строки
    parsingQuery(query, schemaData, clientSocket);
    close(clientSocket);
    cout << "Connection " << clientSocket << " closed" << endl;
    // Освобождение памяти
}

int main() {
    string jsonFileName = "schema.json";
    SchemaInfo schemaData;
    //ввод имени файла и пути
    if (!inputNames(jsonFileName, schemaData)) {
        //close(clientSocket);
        return -1;
    }

    // Создание TCP-сервера
    int serverSocket = socket(AF_INET, SOCK_STREAM, 0); // Создание сокета для прослушивания TCP-соединений
    if (serverSocket == 0) {
        cerr << "Socket creation error" << endl;
        return -1;
    }

    sockaddr_in address; // Структура для хранения адреса сервера
    address.sin_family = AF_INET; // Указание семейства адресов (IPv4)
    address.sin_addr.s_addr = INADDR_ANY; // Привязка к любому IP-адресу, доступному на сервере
    address.sin_port = htons(7432); // Привязка порта 7432 с преобразованием в сетевой порядок байтов

    // Привязка сокета к IP-адресу и порту
    if (bind(serverSocket, (sockaddr*)&address, sizeof(address)) < 0) {
        cerr << "Bind failed" << endl;
        close(serverSocket);
        return -1;
    }

    // Ожидание входящих подключений (до 5 клиентов в очереди).
    if (listen(serverSocket, 5) < 0) {
        cerr << "Listen failed" << endl;
        close(serverSocket);
        return -1;
    }

    cout << "Server is listening on port 7432" << endl;

    vector<thread> clientThreads; // Вектор для хранения потоков, обрабатывающих клиентов.

    while (true) {
        int clientSocket; // Сокет для подключения клиента
        sockaddr_in clientAddress; // Структура для хранения адреса клиента
        socklen_t clientAddressLen = sizeof(clientAddress); // Размер структуры адреса клиента

        // Ожидание подключения клиента
        clientSocket = accept(serverSocket, (sockaddr*)&clientAddress, &clientAddressLen);
        if (clientSocket < 0) {
            cerr << "Failed to accept client" << endl;
            continue;
        }

        cout << "Client " << clientSocket << " connected" << endl;

        // Запуск нового потока для обработки
        clientThreads.emplace_back(thread(handleClient, clientSocket, std::ref(schemaData)));
    }

    // Ждём завершения всех потоков
    for (thread& t : clientThreads) {
        if (t.joinable()) { // Проверка, что поток можно завершить корректно
            t.join();       // Ожидание завершения потока
        }
    }
    
    close(serverSocket); // Закрытие серверного сокета после завершения работы

    return 0;
}