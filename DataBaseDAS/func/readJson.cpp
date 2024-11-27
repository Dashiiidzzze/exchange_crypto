#include "header.h"
#include "../include/json.hpp"

using json = nlohmann::json;

// создание директории
void CreateDir(const string& pathToDir) {
    filesystem::path path(pathToDir);
    if (!filesystem::exists(path)) {
        filesystem::create_directories(path);
    }
}


// создание файла с данными
void CreateFile(const string& pathToFile, const string& fileName, const string& data, bool isDirectory) {
    filesystem::path path(pathToFile);
    if (filesystem::exists(path / fileName)) {
        if (isDirectory) {  // если это файл с таблицей
            ifstream file(path / fileName);
            string line;
            getline(file, line);
            if (line == data) { // данные уже есть в файле
                file.close();
                return;
            }
            file.close();
        } else {
            return;
        }
    }
    // если данные в файле не совпадают с JSON или отсутствуют
    ofstream lockFile(path / fileName);
    if (lockFile.is_open()) {
        lockFile << data;
        lockFile.close();
    } else {
        throw runtime_error("Failed to create lock file in directory");
    }
}


// чтение json файла и создание директорий
void ReadJsonFile(const string& fileName, SchemaInfo& schemaData) {
    ifstream file(schemaData.filepath + "/" + fileName);
    if (!file.is_open()) {
        throw runtime_error("Failed to open schema.json");
    }

    // чтение json
    json schema;
    file >> schema;

    // чтение имени таблицы
    schemaData.name = schema["name"];
    CreateDir(schemaData.name);

    // чтение максимального количества ключей
    schemaData.tuplesLimit = schema["tuples_limit"];

    // чтение структуры таблицы
    json tableStructure = schema["structure"];
    for (auto& [key, value] : tableStructure.items()) {
        // создание директорий
        CreateDir(schemaData.name + "/" + key);
        MyVector<string>* tempValue = CreateVector<string>(10, 50);
        string colNames = key + "_id";
        AddVector(*tempValue, colNames);  // добавлено для чтения индекса
        for (auto columns : value) {
            colNames += ",";
            string temp = columns;
            colNames += temp;
            AddVector(*tempValue, temp);
        }
        CreateFile(schemaData.name + "/" + key, "1.csv", colNames, true);
        //CreateFile(schemaData.name + "/" + key, key + "_lock.txt", "0", false);
        CreateFile(schemaData.name + "/" + key, key + "_pk_sequence.txt", "0", false);
        AddMap<string, MyVector<string>*>(*schemaData.jsonStructure, key, tempValue);
        schemaData.tableMutexes[key];    // Создаем мьютекс для этой таблицы!!!!!!!!!!!!!!!!!!!!!!!!!!!
    }

    file.close();
}