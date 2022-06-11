# Medods-test-task

## Описание
Реализация части сервиса аутентификации. Технологии проекта: MongoDB, Jwt, bcrypt, base64

## Тестирование
Для старта проекта предусмотрены 2 флага addr и dsn, для адреса локального сервера и для соединения с MongoDB
1.  > -dsn="mongodb://localhost:27017/"
2.  > -addr="127.0.0.1:9999"

Команды для быстрого клонирования репозитория и старта проекта
1.  >       git clone https://github.com/Shurubtsov/Medods-test-task.git
2.  >       go run ./cmd/web/

Реализованные ендпоинты:
1. > POST localhost:port/auth/signup?username="имя пользователя"&password="пароль"
2. > POST localhost:port/auth/login?id="ObjectId"
3. > POST localhost:port/auth/refresh

1 - реализует создание пользователя в монгоДБ, параметры запроса username и password

2 - реализует запрос JWT и Refresh токена по айди из монгоДБ, пример "?id=62a34e1b11543aeab2ad021d" , отправляет пару токенов, Refresh токен хешируется base64 и хранится в базе данных в виде bcrypt хеша

3 - реализует запрос на обновление пары JWT и Refresh токенов по отправленному телу из уже имеющегося Refresh токена, который записывается в базу данных при логине

пример тела запроса
>       
>       {
>           "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJBZG1pbiJ9.OqCy2w_NdT4okP_GlKlczPhIsA_H3StuFEj3-ySK1H_iZ-Dlf3s8d9WCq0i_ZP3YcBnTo90wsH4L3Tq8haMEzA",
>           "refresh_token": "YTNlOTdiMTgwM2JjZWVjZDczOGFiMjhhODA2MTBlMjM2ZTRmM2YxYzA4OWJiNWUxMjI5NWZhNzU1YTJkNjE4Mw=="
>       }

`Тестирование API производилось с помощью Postman`
