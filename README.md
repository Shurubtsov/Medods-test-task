# Medods-test-task

## Описание
Реализация части сервиса аутентификации. Технологии проекта: MongoDB-driver, Jwt, bcrypt, base64, mux

## Тестирование
Для старта проекта предусмотрены 2 флага addr и dsn, для адреса локального сервера и для соединения с MongoDB:
-  > -dsn="mongodb://localhost:27017/"
-  > -addr="127.0.0.1:9999"

Команды для быстрого клонирования репозитория и старта проекта
1.  >       git clone https://github.com/Shurubtsov/Medods-test-task.git
2.  >       go run ./cmd/web/


Реализованные ендпоинты:
```go
        // [GET] Маршрут для проверки жизнеспособности вебсервера
	mux.HandleFunc("/", Home(app))
        // [POST] Маршрут для создания пользователя в бд
	mux.HandleFunc("/auth/sign-up", SignUp(app))
    
	// [POST] Основные маршруты для тестового задания
	mux.HandleFunc("/auth/get/tokens", GetTokensForUser(app))
	mux.HandleFunc("/auth/refresh", Refresh(app))
```
Примечание к маршруту `/auth/get/tokens`
> В конце идёт параметр запроса ..tokens?id=GUID вида "cda6498a235d4f7eae19661d41bc154c"

Тело запроса к маршруту `/auth/sign-up`
```json
Пример

{
    "username": "<Имя>",
    "password": "<Пароль>"
}
```

Тело запроса к маршруту `/auth/refresh`
```json
Пример

{
    "access_token": "<Токен>",
    "refresh_token": "<Токен>"
}
```
*Тестирование API производилось с помощью Postman*

### MongoDB

Пример документа внутри бд:
```json
{
  "_id": {
    "$oid": "62a5aa6886ffd3160a6170d1"
  },
  "uuid": "2d02c2395cc749778f2a969ba0629968",
  "username": "Admin",
  "password": "YWRtaW4wMDEwMGFkbWlu",
  "refresh_token": "$2a$14$6g5cZ93iHjX9Ov.E0TzJnu4H1rjx/oMLuOCOkW8flmgQCz5X3OzEa"
}
```
Для работы приложения использовались следующие бд и коллекция:
```go
collection := m.DB.Database("testbase").Collection("users")
```
