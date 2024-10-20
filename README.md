# Поздравлятор

**Поздравлятор** - удобный Telegram бот для поздравления своих коллег, с возможностью напомнить Вам о дне рождении,
написанный на языке Go.

## Функционал
- Регистрация и аутентификация пользователей: пользователи могут регистрироваться и входить в систему.

- Подписка на уведомления о днях рождениях: пользователи могут подписываться на уведомления о днях рождения других пользователей.

- Уведомления через Telegram бота: уведомления о днях рождения доставляются через чат с Telegram ботом.

## Структура проекта

- `cmd/main.go`: точка входа в приложение;

- `internal/auth/auth.go`: модуль для регистрации и аутентификации пользователей;

- `internal/db/database.go`: модуль для подключения к базе данных;

- `internal/notification/bot.go`: модуль для работы с Telegram ботом, включая
обработку команд и взаимодействие с пользователями;

- `internal/db/database.go`: модуль для подключения к базе данных;

- `internal/notification/notification.go`: модуль для управления уведомлениями о днях рождениях;

## Команды

- **/start**: начать взаимодействие с ботом.

- **/signup**: зарегистрировать новый аккаунт.

- **/sublist**: получить список всех пользователей.

- **/sub**: подписаться на уведомления о дне рождения указанного пользователя; 

- **/unsub**: отписаться от уведомлений о дне рождения указанного пользователя.

- **/sublist**: получить список всех пользователей, на которых вы подписаны.

## Используемые технологии
- Telegram API;
- Docker для быстрого развертывания бота;
- Postresql для хранения данных о пользователях. Используется raw sql;
- bcrypt для шифрования пароля пользователя;
- godotenv для хранения переменных БД и ТГ токена;
- golang-migrate для миграции БД.

## Дальнейшая разработка
- Написать unit-тесты;
- Возможность создавать группы из пользователей, которые подписаны на определенного человека,
чтобы они обсудили план поздравления коллеги;
- Использование планировщика cron для автоматизации каких-либо событий.

## Установка и запуск
Склонируйте репозиторий в директорию с помощью `git clone https://github.com/m1hunter/greetings_bot.git`

Перейдите в папку с проектом, и укажите свои переменные окружения в файле `.env`. Также убедитесь, что у вас установлен
`Docker`.

После этого запустите сборку с помощью команды `make`. 

