#include "header.h"

// перезапись во временный файл информации кроме удаленной
void DeleteData(MyVector<string>& tableNames, MyVector<string>& conditionList, SchemaInfo& schemaData) {
    Node* nodeWere = getConditionTree(conditionList);
     
    for (int i = 0; i < tableNames.len; i++) {
        int fileIndex = 1;
        string pathToCSV = schemaData.filepath + "/" + schemaData.name + "/" + tableNames.data[i];
        auto mutexIt = schemaData.tableMutexes.find(tableNames.data[i]);
        if (mutexIt != schemaData.tableMutexes.end()) {
            unique_lock<mutex> lock(mutexIt->second); // Блокировка мьютекса

            while (filesystem::exists(pathToCSV + "/" + to_string(fileIndex) + ".csv")) {
                ifstream file(pathToCSV + "/" + to_string(fileIndex) + ".csv");
                if (!file.is_open()) {
                    throw runtime_error("Failed to open " + (pathToCSV + "/" + to_string(fileIndex) + ".csv"));
                }
                ofstream tempFile(pathToCSV + "/" + to_string(fileIndex) + "_temp.csv");

                string line;
                getline(file, line);
                tempFile << line;
                while (getline(file, line)) {
                    MyVector<string>* row = Split(line, ',');
                    try {
                        if (!isValidRow(nodeWere, *row, *schemaData.jsonStructure, tableNames.data[i])) {
                            tempFile << endl << line;
                        }
                    } catch (const exception& err) {
                        //cerr << err.what() << endl;
                        tempFile.close();
                        file.close();
                        remove((pathToCSV + "/" + to_string(fileIndex) + "_temp.csv").c_str());
                        throw;
                        return;
                    }
                }
                tempFile.close();
                file.close();
                if (remove((pathToCSV + "/" + to_string(fileIndex) + ".csv").c_str()) != 0) {
                    cerr << "Error deleting file" << endl;
                    return;
                }
                if (rename((pathToCSV + "/" + to_string(fileIndex) + "_temp.csv").c_str(), (pathToCSV + "/" + to_string(fileIndex) + ".csv").c_str()) != 0) {
                    cerr << "Error renaming file" << endl;
                    return;
                }

                fileIndex++;
            }
        }
    }
}

// разбиение запроса удаления на кусочки
void ParsingDelete(const MyVector<string>& words, SchemaInfo& schemaData) {
    MyVector<string>* tableNames = CreateVector<string>(5, 50);
    MyVector<string>* conditionList = CreateVector<string>(5, 50);
    int countTabNames = 0;
    int countWereData = 0;
    bool afterWhere = false;
    for (int i = 2; i < words.len; i++ ) {
        if (words.data[i][words.data[i].size() - 1] == ',') {
            words.data[i] = Substr(words.data[i], 0, words.data[i].size() - 1);
        }
        if (words.data[i] == "WHERE") {
            afterWhere = true;
        } else if (afterWhere) {
            AddVector<string>(*conditionList, words.data[i]);
            countWereData++;
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
    if (countTabNames == 0 || countWereData == 0) {
        throw runtime_error("missing table name or data in WERE");
    }

    try {
        DeleteData(*tableNames, *conditionList, schemaData);
    } catch (const exception& err) {
        throw;
        //cerr << err.what()<< endl;
        return;
    }
}