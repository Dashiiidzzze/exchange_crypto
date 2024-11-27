#include "header.h"

// удаление опострафа и проверка синтаксиса
string ApostrovDel(string& str) {
    if (str[str.size() - 1] == ',' && str[str.size() - 2] == ')') {
        str = Substr(str, 0, str.size() - 2);
    } else if (str[str.size() - 1] == ',' || str[str.size() - 1] == ')') {
        str = Substr(str, 0, str.size() - 1);
    }

    if (str[0] == '\'' && str[str.size() - 1] == '\'') {
        str = Substr(str, 1, str.size() - 1);
        return str;
    } else {
        throw runtime_error("invalid sintaxis in VALUES " + str);
    }
}

// проверка количества аргументов относительно столбцов таблиц
void TestAddition(int colLen, const MyVector<string>& tableNames, const MyMap<string, MyVector<string>*>& jsonStructure) {
    for (int i = 0; i < tableNames.len; i++) {
        MyVector<string>* temp = GetMap<string, MyVector<string>*>(jsonStructure, tableNames.data[i]);
        if (temp->len - 1 != colLen) {      // добавить удаление первого элемента из мапа
            throw runtime_error("the number of arguments is not equal to the columns in " + tableNames.data[i]);
        }
    }
}

// чтение файла с количеством записей и перезапись
int PkSequenceRead(const string& path, const bool record, const int newID) {
    fstream pkFile(path);
    if (!pkFile.is_open()) {
        throw runtime_error("Failed to open" + path);
    }
    int lastID = 0;
    if (record) {
        pkFile << newID;
    } else {
        pkFile >> lastID;
    }
    pkFile.close();
    return lastID;
}

// добавление строк в файл
void InsertInTab(MyVector<MyVector<string>*>& addData, MyVector<string>& tableNames, SchemaInfo& schemaData) {
    for (int i = 0; i < tableNames.len; i++) {
        string pathToCSV = schemaData.filepath + "/" + schemaData.name + "/" + tableNames.data[i];
        int lastID = 0;

        // Захватываем мьютекс для таблицы, если она существует в tableMutexes
        auto mutexIt = schemaData.tableMutexes.find(tableNames.data[i]);
        if (mutexIt != schemaData.tableMutexes.end()) {
            unique_lock<mutex> lock(mutexIt->second); // Блокировка мьютекса
            cout << "mutex is locked " << tableNames.data[i] << endl;

            try {
                //BusyTable(pathToCSV, tableNames.data[i] + "_lock.txt", 1);
                lastID = PkSequenceRead(pathToCSV + "/" + tableNames.data[i] + "_pk_sequence.txt", false, 0);
            } catch (const std::exception& err) {
                throw;
                //cerr << err.what() << endl;
                return;
            }

            int newID = lastID;
            for (int j = 0; j < addData.len; j++) {
                newID++;
                string tempPath;
                if (lastID / schemaData.tuplesLimit < newID / schemaData.tuplesLimit) {
                    tempPath = pathToCSV + "/" + to_string(newID / schemaData.tuplesLimit + 1) + ".csv";
                } else {
                    tempPath = pathToCSV + "/" + to_string(lastID / schemaData.tuplesLimit + 1) + ".csv";
                }
                fstream csvFile(tempPath, ios::app);
                if (!csvFile.is_open()) {
                    throw runtime_error("Failed to open" + tempPath);
                }
                csvFile << endl << newID;
                for (int k = 0; k < addData.data[j]->len; k++) {
                    csvFile << "," << addData.data[j]->data[k];
                }
                csvFile.close();
            }
            PkSequenceRead(pathToCSV + "/" + tableNames.data[i] + "_pk_sequence.txt", true, newID);
            //BusyTable(pathToCSV, tableNames.data[i] + "_lock.txt", 0);
        }
    }
}

// разделение запроса вставки на части
void ParsingInsert(const MyVector<string>& words, SchemaInfo& schemaData) {
    MyVector<string>* tableNames = CreateVector<string>(5, 50);
    MyVector<MyVector<string>*>* addData = CreateVector<MyVector<string>*>(10, 50);
    bool afterValues = false;
    int countTabNames = 0;
    int countAddData = 0;
    for (int i = 2; i < words.len; i++) {
        if (words.data[i][words.data[i].size() - 1] == ',') {
            words.data[i] = Substr(words.data[i], 0, words.data[i].size() - 1);
        }
        if (words.data[i] == "VALUES") {
            afterValues = true;
        } else if (afterValues) {
            countAddData++;
            if (words.data[i][0] == '(') {
                MyVector<string>* tempData = CreateVector<string>(5, 50);
                words.data[i] = Substr(words.data[i], 1, words.data[i].size());

                while (words.data[i][words.data[i].size() - 1] != ')' && words.data[i][words.data[i].size() - 2] != ')') {
                    try {
                        ApostrovDel(words.data[i]);
                    } catch (const exception& err) {
                        throw;
                        //cerr << err.what() << words.data[i] << endl;
                        return;
                    }
                    
                    AddVector<string>(*tempData, words.data[i]);
                    i++;
                }
                try {
                    ApostrovDel(words.data[i]);
                    AddVector<string>(*tempData, words.data[i]);
                    TestAddition(tempData->len, *tableNames, *schemaData.jsonStructure);
                } catch (const exception& err) {
                    throw;
                    //cerr << err.what() << endl;
                    return;
                }
                AddVector<MyVector<string>*>(*addData, tempData);
            }
            
        } else {
            countTabNames++;
            try {
                GetMap(*schemaData.jsonStructure, words.data[i]);
            } catch (const exception& err) {
                throw;
                //cerr << err.what() << ": table " << words.data[i] << " is missing" << endl;
                return;
            }
            AddVector<string>(*tableNames, words.data[i]);
        }
    }
    if (countTabNames == 0 || countAddData == 0) {
        throw runtime_error("missing table name or data in VALUES");
    }

    try {
        InsertInTab(*addData, *tableNames, schemaData);
    } catch (const exception& err) {
        throw;
        //cerr << err.what() << endl;
        return;
    }
}