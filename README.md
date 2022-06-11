# Medods-test-task

## Описание
Реализация части сервиса аутентификации. Технологии проекта: MongoDB, Jwt, bcrypt, base64

## Тестирование
Для старта проекта предусмотрены 2 флага addr и dsn, для адреса локального сервера и для соединения с MongoDB:
-  > -dsn="mongodb://localhost:27017/"
-  > -addr="127.0.0.1:9999"

Команды для быстрого клонирования репозитория и старта проекта
1.  >       git clone https://github.com/Shurubtsov/Medods-test-task.git
2.  >       go run ./cmd/web/


Реализованные ендпоинты:
1. > POST localhost:port/auth/sign-up
2. > POST localhost:port/auth/login?id="ObjectId"
3. > POST localhost:port/auth/refresh

1 - реализует создание пользователя в монгоДБ, параметры запроса username и password

```json
Пример запроса.

{
    "username": "<Имя>",
    "password": "<Пароль>"
}
```

2 - реализует запрос JWT и Refresh токена по айди из монгоДБ, пример "?id=62a34e1b11543aeab2ad021d" , отправляет пару токенов, Refresh токен хешируется base64 и хранится в базе данных в виде bcrypt хеша

3 - реализует запрос на обновление пары JWT и Refresh токенов по отправленному телу из уже имеющегося Refresh токена, который записывается в базу данных при логине

```json
Пример запроса.

{
    "access_token": "<Токен>",
    "refresh_token": "<Токен>"
}
```
`Тестирование API производилось с помощью Postman`
