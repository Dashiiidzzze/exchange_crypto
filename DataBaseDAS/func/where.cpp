#include "header.h"

string ApostDel(string str) {
    if (str[0] == '\'' && str[str.size() - 1] == '\'') {
        str = Substr(str, 1, str.size() - 1);
        return str;
    } else {
        throw runtime_error("invalid sintaxis in WHERE " + str);
    }
}

// Вспомогательная функция для разделения строки по оператору
MyVector<MyVector<string>*>* splitByOperator(const MyVector<string>& query, const string& op) {
    MyVector<string>* left = CreateVector<string>(6, 50);
    MyVector<string>* right = CreateVector<string>(6, 50);
    bool afterOp = false;
    for (int i = 0; i < query.len; i++) {
        if (query.data[i] == op) {
            afterOp = true;
        } else if (afterOp) {
            AddVector(*right, query.data[i]);
        } else {
            AddVector(*left, query.data[i]);
        }
    }
    MyVector<MyVector<string>*>* parseVector = CreateVector<MyVector<string>*>(5, 50);
    if (afterOp) {
        AddVector(*parseVector, left);
        AddVector(*parseVector, right);
        
    } else {
        AddVector(*parseVector, left);
    }
    return parseVector;
}


Node* getConditionTree(const MyVector<string>& query) {
    MyVector<MyVector<string>*>* orParts = splitByOperator(query, "OR");

    // OR
    if (orParts->len > 1) {
        Node* root = new Node(NodeType::OrNode);
        root->left = getConditionTree(*orParts->data[0]);
        root->right = getConditionTree(*orParts->data[1]);
        return root;
    }
    // AND
    MyVector<MyVector<std::string>*>* andParts = splitByOperator(query, "AND");
    if (andParts->len > 1) {
        Node* root = new Node(NodeType::AndNode);
        root->left = getConditionTree(*andParts->data[0]);
        root->right = getConditionTree(*andParts->data[1]);
        return root;
    }

    // Simple condition
    return new Node(NodeType::ConditionNode, query);
}

bool isValidRow(Node* node, const MyVector<string>& row, const MyMap<string, MyVector<string>*>& jsonStructure, const string& tabNames) {
    if (!node) {
        return false;
    }

    switch (node->nodeType) {
    case NodeType::ConditionNode: {
        if (node->value.len != 3) {
            return false;
        }

        MyVector<string> *part1Splitted = Split(node->value.data[0], '.');
        if (part1Splitted->len != 2) {
            return false;
        }
    
        // существует ли запрашиваемая таблица
        int columnIndex = -1;
        try {
            MyVector<string>* colNames = GetMap(jsonStructure, part1Splitted->data[0]);
            for (int i = 0; i < colNames->len; i++) {
                if (colNames->data[i] == part1Splitted->data[1]) {
                    columnIndex = i;
                    break;
                }
            }
        } catch (const exception& err) {
            throw;
            //cerr << err.what() << ": table " << part1Splitted->data[0] << " is missing" << std::endl;
            return false;
        }

        if (columnIndex == -1) {
            cerr << "Column " << part1Splitted->data[1] << " is missing in table " << part1Splitted->data[0] << std::endl;
            return false;
        }

        string delApostr = ApostDel(node->value.data[2]);
        //if (tabNames == part1Splitted->data[0] && row.data[columnIndex + 1] == delApostr) {       // изменить здесь
        if (tabNames == part1Splitted->data[0] && row.data[columnIndex] == delApostr) {  
            return true;
        }

        return false;
    }
    case NodeType::OrNode:
        return isValidRow(node->left, row, jsonStructure, tabNames) ||
                isValidRow(node->right, row, jsonStructure, tabNames);
    case NodeType::AndNode:
        return isValidRow(node->left, row, jsonStructure, tabNames) &&
                isValidRow(node->right, row, jsonStructure, tabNames);
    default:
        return false;
    }
}