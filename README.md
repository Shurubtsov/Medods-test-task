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
> В конце идёт параметр запроса ..tokens?id=<GUID>

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
