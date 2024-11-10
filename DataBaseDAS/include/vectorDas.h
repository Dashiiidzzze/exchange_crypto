#ifndef VECTORDAS_H
#define VECTORDAS_H

#include <iostream>
#include <iomanip>


template <typename T>
struct MyVector {
    T* data;      //сам массив
    size_t len;        //длина
    size_t cap;        //capacity - объем
    size_t LoadFactor; //с какого процента заполнения увеличиваем объем = 50%
};


// перегрузка оператора вывода
template <typename T>
std::ostream& operator << (std::ostream& os, const MyVector<T>& vec) {
    for (size_t i = 0; i < vec.len; i++) {
        std::cout << vec.data[i];
        if (i < vec.len - 1) std::cout << std::setw(25);
    }
    return os;
}

// инициализация вектора
template <typename T>
MyVector<T>* CreateVector(size_t initCapacity, size_t initLoadFactor) {
    if (initCapacity <= 0 || initLoadFactor <= 0 || initLoadFactor > 100) {
        throw std::runtime_error("Index out of range");
    }

    MyVector<T>* vec = new MyVector<T>;
    vec->data = new T[initCapacity];
    vec->len = 0;
    vec->cap = initCapacity;
    vec->LoadFactor = initLoadFactor;
    return vec;
}

// увеличение массива
template <typename T>
void Expansion(MyVector<T>& vec) {
    size_t newCap = vec.cap * 2;
    T* newData = new T[newCap];
    for (size_t i = 0; i < vec.len; i++) {     //копируем данные из старого массива в новый
        newData[i] = vec.data[i];
    }
    delete[] vec.data;                      // очистка памяти
    vec.data = newData;
    vec.cap = newCap;
}

// добавление элемента в вектор
template <typename T>
void AddVector(MyVector<T>& vec, T value) {
    if ((vec.len + 1) * 100 / vec.cap >= vec.LoadFactor) { //обновление размера массива
        Expansion(vec);
    }
    vec.data[vec.len] = value;
    vec.len++;
}


//удаление элемента из вектора
template <typename T>
void DeleteVector(MyVector<T>& vec, size_t index) {
    if (index < 0 || index >= vec.len) {
        throw std::runtime_error("Index out of range");
    }

    for (size_t i = index; i < vec.len - 1; i++) {
        vec.data[i] = vec.data[i + 1];
    }

    vec.len--;
}


// замена элемента по индексу
template <typename T>
void ReplaceVector(MyVector<T>& vec, size_t index, T value) {
    if (index < 0 || index >= vec.len) {
        throw std::runtime_error("Index out of range");
    }
    vec.data[index] = value;
}

#endif