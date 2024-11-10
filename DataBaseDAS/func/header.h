#pragma once

#include <iostream>
#include <fstream>
#include <filesystem>
#include <arpa/inet.h>
#include <unistd.h>

#include <map>
#include <mutex>

#include "../include/mapDas.h"
#include "../include/vectorDas.h"

using namespace std;

struct SchemaInfo {
    string filepath = ".";
    string name;
    int tuplesLimit;
    MyMap<string, MyVector<string>*>* jsonStructure;
    map<string, mutex> tableMutexes;
};

// support functions
int Strlen(const string &str);
string Substr(const string &str, int start, int end);
MyVector<string>* Split(const string &str, char delim);


// reading json file
void CreateDir(const string& pathToDir);
void CreateFile(const string& pathToFile, const string& fileName, const string& data, bool isDirectory);
void ReadJsonFile(const string& fileName, SchemaInfo& schemaData);

// where
// Тип узла
enum class NodeType {
    ConditionNode,
    OrNode,
    AndNode
};

// Структура
struct Node {
    NodeType nodeType;
    MyVector<std::string> value;
    Node* left;
    Node* right;

    Node(NodeType type, const MyVector<std::string> val = {}, Node* l = nullptr, Node* r = nullptr)
        : nodeType(type), value(val), left(l), right(r) {}
};

string ApostDel(string str);
MyVector<MyVector<string>*>* splitByOperator(const MyVector<string>& query, const string& op);
Node* getConditionTree(const MyVector<string>& query);
bool isValidRow(Node* node, const MyVector<string>& row, const MyMap<string, MyVector<string>*>& jsonStructure, const string& tabNames);



// select
bool AllVritingToVec(Node* nodeWere, const string& tableName, string& line, MyVector<MyVector<string>*>& tabData, SchemaInfo& schemaData, bool where);
bool VritingToVec(Node* nodeWere, const string& tableName, string& line, MyVector<MyVector<string>*>& tabData, SchemaInfo& schemaData, bool where, MyVector<int>& colIndex);
MyVector<MyVector<string>*>* ReadTable(const string& tableName, SchemaInfo& schemaData, const MyVector<string>& colNames, const MyVector<string>& conditionList, bool where);
void DecartMult(const MyVector<MyVector<MyVector<string>*>*>& tablesData, MyVector<MyVector<string>*>& temp, int counterTab, int tab, int clientSocket);
void PreparationSelect(const MyVector<string>& colNames, const MyVector<string>& tableNames, const MyVector<string>& conditionList, SchemaInfo& schemaData, bool where, int clientSocket);
void ParsingSelect(const MyVector<string>& words, SchemaInfo& schemaData, int clientSocket);


// insert
string ApostrovDel(string& str);
void TestAddition(int colLen, const MyVector<string>& tableNames, const MyMap<string, MyVector<string>*>& jsonStructure);
int PkSequenceRead(const string& path, const bool record, const int newID);
void InsertInTab(MyVector<MyVector<string>*>& addData, MyVector<string>& tableNames, SchemaInfo& schemaData);
void ParsingInsert(const MyVector<string>& words, SchemaInfo& schemaData);

// delete
void DeleteData(MyVector<string>& tableNames, MyVector<string>& conditionList, SchemaInfo& schemaData);
void ParsingDelete(const MyVector<string>& words, SchemaInfo& schemaData);