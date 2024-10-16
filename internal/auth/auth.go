package auth

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"projectik/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthBot struct {
	api *tgbotapi.BotAPI
	db  *sql.DB
}

func NewAuthBot(api *tgbotapi.BotAPI, db *sql.DB) *AuthBot {
	return &AuthBot{
		api: api,
		db:  db,
	}
}

func (b *AuthBot) HandleMessage(update tgbotapi.Update) {
	if update.Message == nil || !update.Message.IsCommand() {
		return
	}

	command := update.Message.Command()

	switch command {
	case "signup":
		b.handleSignup(update.Message, update)
	case "whoami":
		b.handleLogin(update.Message) // теперь "whoami" обрабатывается в handleLogin
	default:
		b.handleUnknownCommand(update.Message)
	}
}

func (b *AuthBot) handleSignup(msg *tgbotapi.Message, update tgbotapi.Update) {
	args := msg.CommandArguments()
	if args == "" {
		reply := "Для регистрации введите /signup <password> <firstname> <lastname> <birthday>"
		b.sendMessage(msg.Chat.ID, reply)
		return
	}

	params := strings.Fields(args)

	if len(params) < 4 {
		reply := "Неверный формат. Используйте /signup <password> <firstname> <lastname> <birthday>"
		b.sendMessage(msg.Chat.ID, reply)
		return
	}

	username := update.Message.From.UserName
	if username == "" {
		reply := "У вас отсутствует никнейм. Чтобы использовать бота, укажите в настройках никнейм."
		b.sendMessage(msg.Chat.ID, reply)
		return
	}
	password := params[0]
	firstname := params[1]
	lastname := params[2]
	birthdayStr := params[3]

	// Проверяем, существует ли пользователь с данным username
	var existingUser models.User
	err := b.db.QueryRow("SELECT username FROM users WHERE username = $1", username).Scan(&existingUser.Username)
	if err == nil {
		reply := "Вы уже зарегистрированы, ваш Telegram ID: " + msg.From.UserName
		b.sendMessage(msg.Chat.ID, reply)
		return
	}

	birthday, err := time.Parse("2006-01-02", birthdayStr)
	if err != nil {
		reply := "Неверная дата рождения. Убедитесь, что Вы ввели дату в формате 2006-01-02"
		b.sendMessage(msg.Chat.ID, reply)
		return
	}
	//	вообще, в текущей реализации пароль не нужен, но он сделан просто чтобы быть, возможно
	// 	в будущем будет как минимум basic auth
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		reply := "Ошибка при регистрации. Повторите попытку позже."
		b.sendMessage(msg.Chat.ID, reply)
		log.Printf("User %s registration error: ", msg.From.UserName, err)
		return
	}

	chatId := update.Message.Chat.ID

	// добавляем нового пользователя в базу данных
	_, err = b.db.Exec(
		"INSERT INTO users (username, password, firstname, lastname, birthday, chat_id) VALUES ($1, $2, $3, $4, $5, $6)",
		username, hashedPassword, firstname, lastname, birthday, chatId)

	if err != nil {
		reply := "Ошибка при регистрации пользователя: " + err.Error()
		b.sendMessage(msg.Chat.ID, reply)
		log.Printf("Error during user registration: %v", err)
		return
	}

	reply := "Пользователь зарегистрирован успешно!"
	b.sendMessage(msg.Chat.ID, reply)
	log.Printf("User %s registered successfully", username)
}

func (b *AuthBot) handleLogin(msg *tgbotapi.Message) {
	// Команда whoami выводит информацию о пользователе
	var user models.User
	err := b.db.QueryRow("SELECT username FROM users WHERE chat_id = $1", msg.Chat.ID).Scan(&user.Username)
	if err != nil {
		reply := "Вы не авторизованы. Пожалуйста, зарегистрируйтесь с помощью /signup."
		b.sendMessage(msg.Chat.ID, reply)
		return
	}

	reply := "Ваш никнейм в системе: " + user.Username
	b.sendMessage(msg.Chat.ID, reply)
}

func (b *AuthBot) handleUnknownCommand(msg *tgbotapi.Message) {
	reply := "Неизвестная команда. Используйте /start для начала работы или /help для получения списка команд."
	b.sendMessage(msg.Chat.ID, reply)
}

func (b *AuthBot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}
