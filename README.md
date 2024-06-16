# Поздравлятор

**Поздравлятор** - удобный Telegram бот для поздравления своих коллег, с возможностью напомнить Вам о дне рождении.
Функционал бота:
- Регистрация и аутентификация пользователей: пользователи могут регистрироваться и входить в систему.

- Подписка на уведомления о днях рождениях: пользователи могут подписываться на уведомления о днях рождения других пользователей.

- Уведомления через Telegram бота: уведомления о днях рождения доставляются через чат с Telegram ботом.

Бот использует ЯП Go и базу данных PostgreSQL.
Для запуска бота необходимо создать файл окружения .env в корне проекта и заполнить его следующим образом:
```
TELEGRAM_BOT_TOKEN=TgBot_Token
DB_HOST=db_host
DB_PORT=db_port
DB_USER=db_user
DB_PASSWORD=db_pass
DB_NAME=db_name
```
Структура проекта

_cmd/app/main.go: точка входа в приложение.
internal/auth/auth.go: модуль для регистрации и аутентификации пользователей.
internal/bot/bot.go: модуль для работы с Telegram ботом, включая обработку команд и взаимодействие с пользователями.
internal/db/db.go: модуль для подключения к базе данных и выполнения запросов.
internal/notification/notification.go: модуль для управления уведомлениями, использующий cron для планирования задач._

Команды Бота

_/start: начать взаимодействие с ботом.
/signup: зарегистрировать новый аккаунт.
/userslist: получить список всех пользователей.
/sub: подписаться на уведомления о дне рождения указанного пользователя.
/unsub: отписаться от уведомлений о дне рождения указанного пользователя.
/list: получить список всех пользователей.
/sublist: получить список всех пользователей, на которых вы подписаны._