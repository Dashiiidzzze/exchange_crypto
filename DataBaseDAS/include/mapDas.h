#ifndef MAPDAS_H
#define MAPDAS_H

#include <iostream>
#include <string>
//#include "vectorDas.h"

//using namespace std;

// структура для хранения значения
template <typename TK, typename TV>
struct Nod {
    TK key;
    TV value;
    Nod* next;
};

// структура для хранения ключа и значения
template <typename TK, typename TV>
struct MyMap {
    Nod<TK, TV>** data;
    size_t len;
    size_t cap;
    size_t LoadFactor;
};

// хэш-функция для ключа string
template <typename TK>
int HashCode(const TK& key, const int capacity) {
    unsigned long hash = 5381;
    int c = 0;
    for (char ch : key) {
        hash = ((hash << 5) + hash) + ch; /* hash * 33 + c */
    }
    return hash % capacity;
}


// инициализация хэш таблицы
template <typename TK, typename TV>
MyMap<TK, TV>* CreateMap(int initCapacity, int initLoadFactor) {
    if (initCapacity <= 0 || initLoadFactor <= 0 || initLoadFactor > 100) {
        throw std::runtime_error("Index out of range");
    }

    MyMap<TK, TV>* map = new MyMap<TK, TV>;
    map->data = new Nod<TK, TV>*[initCapacity];
    for (size_t i = 0; i < initCapacity; i++) {
        map->data[i] = nullptr;
    }

    map->len = 0;
    map->cap = initCapacity;
    map->LoadFactor = initLoadFactor;
    return map;
}

// расширение
template <typename TK, typename TV>
void Expansion(MyMap<TK, TV>& map) {
    size_t newCap = map.cap * 2;
    Nod<TK, TV>** newData = new Nod<TK, TV>*[newCap];
    for (size_t i = 0; i < newCap; i++) {
        newData[i] = nullptr;
    }
    // проход по всем ячейкам
    for (size_t i = 0; i < map.cap; i++) {
        Nod<TK, TV>* curr = map.data[i];
        // проход по парам коллизионных значений и обновление
        while (curr != nullptr) {
            Nod<TK, TV>* next = curr->next;
            size_t index = HashCode(curr->key, newCap);
            curr->next = newData[index];
            newData[index] = curr;
            curr = next;
        }
    }

    delete[] map.data;

    map.data = newData;
    map.cap = newCap;
}

// обработка коллизий
template <typename TK, typename TV>
void CollisionManage(MyMap<TK, TV>& map, int index, const TK& key, const TV& value) {
    Nod<TK, TV>* newNod = new Nod<TK, TV>{key, value, nullptr};
    Nod<TK, TV>* curr = map.data[index];
    while (curr->next != nullptr) {
        curr = curr->next;
    }
    curr->next = newNod;
}

// добавление элементов
template <typename TK, typename TV>
void AddMap(MyMap<TK, TV>& map, const TK& key, const TV& value) {
    if ((map.len + 1) * 100 / map.cap >= map.LoadFactor) {
        Expansion(map);
    }
    size_t index = HashCode(key, map.cap);
    Nod<TK, TV>* temp = map.data[index];
    if (temp != nullptr) {
        while (temp != nullptr) {
            if (temp->key == key) {
                // Элемент уже существует, обновить значение
                temp->value = value;
                map.data[index] = temp;
                return;
            }
            temp = temp->next;
        }
        CollisionManage(map, index, key, value);
    } else {
        Nod<TK, TV>* newNod = new Nod<TK, TV>{key, value, map.data[index]};
        map.data[index] = newNod;
        map.len++;
    }

}


// поиск элементов по ключу
template <typename TK, typename TV>
TV GetMap(const MyMap<TK, TV>& map, const TK& key) {
    size_t index = HashCode(key, map.cap);
    Nod<TK, TV>* curr = map.data[index];
    while (curr != nullptr) {
        if (curr->key == key) {
            return curr->value;
        }
        curr = curr->next;
    }

    throw std::runtime_error("Key not found");
}


// удаление элементов
template <typename TK, typename TV>
void DeleteMap(MyMap<TK, TV>& map, const TK& key) {
    size_t index = HashCode(key, map.cap);
    Nod<TK, TV>* curr = map.data[index];
    Nod<TK, TV>* prev = nullptr;
    while (curr != nullptr) {
        if (curr->key == key) {
            if (prev == nullptr) {
                map.data[index] = curr->next;
            } else {
                prev->next = curr->next;
            }
            delete curr;
            map.len--;
            return;
        }
        prev = curr;
        curr = curr->next;
    }
    throw std::runtime_error("Key not found");
}


// очистка памяти
template <typename TK, typename TV>
void DestroyMap(MyMap<TK, TV>& map) {
    for (size_t i = 0; i < map.cap; i++) {
        Nod<TK, TV>* curr = map.data[i];
        while (curr != nullptr) {
            Nod<TK, TV>* next = curr->next;
            delete curr;
            curr = next;
        }
    }
    delete[] map.data;
    map.data = nullptr;
    map.len = 0;
    map.cap = 0;
}


#endif
