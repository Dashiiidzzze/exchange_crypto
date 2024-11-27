## Валютная биржа

для запуска сервера на свооем пк:

1) `git clone https://github.com/Dashiiidzzze/exchange_crypto/tree/main`
2) `docker compose up --build`
3) начать проводить HTTP запросы (через браузер или утилиту curl)

### Примеры HTTP запросов

Создание пользователей:
`curl -X POST http://localhost:8080/user -H "Content-Type: application/json" -d '{"username": "dasha"}'`

Создание ордера:

```bash
curl -X POST http://localhost:8080/order \
     -H "Content-Type: application/json" \
     -H "X-USER-KEY: qwertykey" \
     -d '{
           "pair_id": 1,
           "quantity": 100.5,
           "price": 2500.75,
           "type": "buy"
         }'
```

Получение списка ордеров:
`curl -X GET http://localhost:8080/order -H "Content-Type: application/json"`

Удаление ордера:

```bash
curl -X DELETE http://localhost:8080/order \
     -H "Content-Type: application/json" \
     -H "X-USER-KEY: qwertykey" \
     -d '{
           "order_id": 2
         }'
```

Получение информации о лотах:
`curl -X GET http://localhost:8080/lot -H "Content-Type: application/json"`

Получение информации о парах:
`curl -X GET http://localhost:8080/pair -H "Content-Type: application/json"`

Получение информации об активах пользователя:
`curl -X GET http://localhost:8080/balance -H "Content-Type: application/json" -H "X-USER-KEY: qwertykey"`
