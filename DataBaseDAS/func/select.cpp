#include "header.h"

// считвыание всех подходящих строк 
bool AllVritingToVec(Node* nodeWere, const string& tableName, string& line, MyVector<MyVector<string>*>& tabData, SchemaInfo& schemaData, bool where) {
    MyVector<string>* row = Split(line, ',');
    if (where) {
        try {
            if (isValidRow(nodeWere, *row, *schemaData.jsonStructure, tableName)) {
                //DeleteVector<string>(*row, 0);          // удалить
                AddVector(tabData, row);
            }
        } catch (const exception& err) {
            throw;
            //cerr << err.what() << endl;
            return false;
        }
    } else {
        //DeleteVector<string>(*row, 0);              // удалить
        AddVector(tabData, row);
    }
    return true;
}

// считывание подходящих строк из выбранных столбцов
bool VritingToVec(Node* nodeWere, const string& tableName, string& line, MyVector<MyVector<string>*>& tabData, SchemaInfo& schemaData, bool where, MyVector<int>& colIndex) {
    MyVector<string>* row = Split(line, ',');
    MyVector<string>* newRow = CreateVector<string>(colIndex.len, 50);
    if (where) {
        try {
            if (isValidRow(nodeWere, *row, *schemaData.jsonStructure, tableName)) {
                for (int i = 0; i < colIndex.len; i++) {
                    AddVector(*newRow, row->data[colIndex.data[i]]);
                }
                AddVector(tabData, newRow);
            }
        } catch (const exception& err) {
            throw;
            //cerr << err.what() << endl;
            return false;
        }
    } else {
        for (int i = 0; i < colIndex.len; i++) {
            AddVector(*newRow, row->data[colIndex.data[i]]);
        }
        AddVector(tabData, newRow);
    }
    return true;
}



// чтение таблицы из файла
MyVector<MyVector<string>*>* ReadTable(const string& tableName, SchemaInfo& schemaData, const MyVector<string>& colNames, const MyVector<string>& conditionList, bool where) {
    MyVector<MyVector<string>*>* tabData = CreateVector<MyVector<string>*>(5, 50);
    string pathToCSV = schemaData.filepath + "/" + schemaData.name + "/" + tableName;
    int fileIndex = 1;

    // Захватываем мьютекс для таблицы, если она существует в tableMutexes
    auto mutexIt = schemaData.tableMutexes.find(tableName);
    if (mutexIt != schemaData.tableMutexes.end()) {
        unique_lock<mutex> lock(mutexIt->second); // Блокировка мьютекса
        
        Node* nodeWere = getConditionTree(conditionList);
        while (filesystem::exists(pathToCSV + "/" + to_string(fileIndex) + ".csv")) {
            ifstream file(pathToCSV + "/" + to_string(fileIndex) + ".csv");
            if (!file.is_open()) {
                throw runtime_error("Failed to open file" + (pathToCSV + "/" + to_string(fileIndex) + ".csv"));
            }
            string firstLine;
            getline(file, firstLine);
            if (colNames.data[0] == "*") {
                string line;
                while (getline(file, line)) {
                    if (!AllVritingToVec(nodeWere, tableName, line, *tabData, schemaData, where)) {
                        file.close();
                        return tabData;
                    }
                }
            } else {
                MyVector<string>* fileColNames = GetMap<string, MyVector<string>*>(*schemaData.jsonStructure, tableName);
                MyVector<int>* colIndex = CreateVector<int>(10, 50);
                for (int i = 0; i < fileColNames->len; i++) {
                    for (int j = 1; j < colNames.len; j++) {
                        if (fileColNames->data[i] == colNames.data[j]) {
                            //AddVector(*colIndex, i + 1);                // убрать +1
                            AddVector(*colIndex, i);
                        }
                    }
                }
                string line;
                while (getline(file, line)) {
                    if (!VritingToVec(nodeWere, tableName, line, *tabData, schemaData, where, *colIndex)) {
                        file.close();
                        return tabData;
                    }
                }
            }

            file.close();
            fileIndex += 1;
        }
    }
    return tabData;
}


// вывод содержимого таблиц в виде декартового произведения
void DecartMult(const MyVector<MyVector<MyVector<string>*>*>& tablesData, MyVector<MyVector<string>*>& temp, int counterTab, int tab, int clientSocket) {
    for (int i = 0; i < tablesData.data[counterTab]->len; i++) {
        temp.data[counterTab] = tablesData.data[counterTab]->data[i];

        if (counterTab < tab - 1) {
            DecartMult(tablesData, temp, counterTab + 1, tab, clientSocket);
        } else {
            for (int j = 0; j < tab; j++) {
                for (int k = 0; k < temp.data[j]->len; k++) {
                    send(clientSocket, (temp.data[j]->data[k] + " ").c_str(), (temp.data[j]->data[k] + " ").size(), 0);
                }
                //cout << *temp.data[j] << setw(25);
            }
            string enter = "\n";
            send(clientSocket, enter.c_str(), enter.size(), 0);
            //cout << endl;
        }
    }

    return;
}

// подготовка к чтению и выводу данных
void PreparationSelect(const MyVector<string>& colNames, const MyVector<string>& tableNames, const MyVector<string>& conditionList, SchemaInfo& schemaData, bool where, int clientSocket) {
    MyVector<MyVector<MyVector<string>*>*>* tablesData = CreateVector<MyVector<MyVector<string>*>*>(10, 50);
    if (colNames.data[0] == "*") {      // чтение всех данных из таблиц
        for (int j = 0; j < tableNames.len; j++) {
            MyVector<MyVector<string>*>* tableData = ReadTable(tableNames.data[j], schemaData, colNames, conditionList, where);
            AddVector(*tablesData, tableData);
        }
    } else {
        for (int i = 0; i < tableNames.len; i++) {
            MyVector<string>* tabColPair = CreateVector<string>(5, 50);
            AddVector(*tabColPair, tableNames.data[i]);
            for (int j = 0; j < colNames.len; j++) {
                MyVector<string>* splitColNames = Split(colNames.data[j], '.');
                try {
                    GetMap(*schemaData.jsonStructure, splitColNames->data[0]);
                } catch (const exception& err) {
                    throw;
                    //cerr << err.what() << ": table " << splitColNames->data[0] << " is missing" << endl;
                    return;
                }
                if (splitColNames->data[0] == tableNames.data[i]) {
                    AddVector(*tabColPair, splitColNames->data[1]);
                }
            }
            MyVector<MyVector<string>*>* tableData = ReadTable(tabColPair->data[0], schemaData, *tabColPair, conditionList, where);
            AddVector(*tablesData, tableData);
        }
    }

    MyVector<MyVector<string>*>* temp = CreateVector<MyVector<string>*>(tablesData->len * 2, 50);
    string resStr;
    DecartMult(*tablesData, *temp, 0, tablesData->len, clientSocket);
    return;
}

// парсинг SELECT запроса
void ParsingSelect(const MyVector<string>& words, SchemaInfo& schemaData, int clientSocket) {
    MyVector<string>* colNames = CreateVector<string>(10, 50);          // названия колонок в формате таблица1.колонка1
    MyVector<string>* tableNames = CreateVector<string>(10, 50);        // названия таблиц в формате  таблица1
    MyVector<string>* conditionList = CreateVector<string>(10, 50);     // список условий where
    bool afterFrom = false;
    bool afterWhere = false;
    int countTabNames = 0;
    int countData = 0;
    int countWhereData = 0;
    for (int i = 1; i < words.len; i++) {
        if (words.data[i][words.data[i].size() - 1] == ',') {
            words.data[i] = Substr(words.data[i], 0, words.data[i].size() - 1);
        }
        if (words.data[i] == "FROM") {
            afterFrom = true;
        } else if (words.data[i] == "WHERE") {
            afterWhere = true;
        } else if (afterWhere) {
            countWhereData++;
            AddVector<string>(*conditionList, words.data[i]);
        } else if (afterFrom) {
            try {
                GetMap(*schemaData.jsonStructure, words.data[i]);
            } catch (const exception& err) {
                throw;
                //cerr << err.what() << ": table " << words.data[i] << " is missing" << endl;
                return;
            }
            countTabNames++;
            AddVector(*tableNames, words.data[i]);
        } else {
            countData++;
            AddVector(*colNames, words.data[i]);
        }
    }
    if (countTabNames == 0 || countData == 0) {
        throw runtime_error("missing table name or data in FROM");
    }
    if (countWhereData == 0) { //const string& schemaName, const string& filePath, const MyMap<string, MyVector<string>*>& jsonStructure
        PreparationSelect(*colNames, *tableNames, *conditionList, schemaData, false, clientSocket);
    } else {
        PreparationSelect(*colNames, *tableNames, *conditionList, schemaData, true, clientSocket);
    }
}