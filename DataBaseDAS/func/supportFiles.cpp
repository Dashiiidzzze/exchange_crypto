#include "header.h"

// Возвращает длину строки
int Strlen(const string &str) {
    int len = 0;
    while (str[len]!= '\0') {
        len++;
    }
    return len;
}

// Возвращает подстроку от start до end (не включая end)
string Substr(const string &str, int start, int end) {
    string result;
    for (int i = start; i < end; i++) {
        result += str[i];
    }
    return result;
}

// Разбивает строку на слова с разделителем delim.
MyVector<string>* Split(const string &str, char delim) {
    int index = 0;
    MyVector<string>* words = CreateVector<string>(10, 50);
    int length = Strlen(str);

    while (true) {
        int delimIndex = index;
        while (str[delimIndex] != delim && delimIndex != length) delimIndex++;

        string word = Substr(str, index, delimIndex);
        AddVector(*words, word);
        index = delimIndex + 1;
        if (delimIndex == length) break;
    }

    return words;
}

// проверка на занятость таблицы длугим пользоватклкм
// void BusyTable(const string& path, const string& fileName, const int rank) {
//     fstream lockFile(path + "/" + fileName);
//     if (!lockFile.is_open()) {
//         throw runtime_error("Failed to open " + (path + "/" + fileName));
//     }
//     int lock = 0;
//     lockFile >> lock;
//     if (lock == 1 && rank == 1) {
//         lockFile.close();
//         throw runtime_error("Table " + fileName + " is locked by another process");
//     } else {
//         lockFile << rank;
//     }
//     lockFile.close();
// }