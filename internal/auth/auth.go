package auth

import (
	"log"
	"strings"

	"projectik/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthBot struct {
	api *tgbotapi.BotAPI
	db  *gorm.DB
}

func NewAuthBot(api *tgbotapi.BotAPI, db *gorm.DB) *AuthBot {
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
		b.handleSignup(update.Message)
	case "login":
		b.handleLogin(update.Message)
	default:
		b.handleUnknownCommand(update.Message)
	}
}

func (b *AuthBot) handleSignup(msg *tgbotapi.Message) {
	args := msg.CommandArguments()
	if args == "" {
		reply := "Для регистрации введи /signup <username> <password> <firstname> <lastname> <birthday>"
		b.sendMessage(msg.Chat.ID, reply)
		return
	}

	params := strings.Fields(args)
	if len(params) < 5 {
		reply := "Неверный формат. Используй /signup <username> <password> <firstname> <lastname> <birthday>"
		b.sendMessage(msg.Chat.ID, reply)
		return
	}

	username := params[0]
	password := params[1]
	firstname := params[2]
	lastname := params[3]
	birthday := params[4]

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		reply := "Ошибка при шифровании пароля: " + err.Error()
		b.sendMessage(msg.Chat.ID, reply)
		return
	}

	user := models.User{
		Username:  username,
		Password:  string(hashedPassword),
		FirstName: firstname,
		LastName:  lastname,
		Birthday:  birthday,
	}

	if err := b.db.Create(&user).Error; err != nil {
		reply := "Ошибка при регистрации пользователя: " + err.Error()
		b.sendMessage(msg.Chat.ID, reply)
		return
	}

	reply := "Пользователь зарегистрирован успешно!"
	b.sendMessage(msg.Chat.ID, reply)
}

func (b *AuthBot) handleLogin(msg *tgbotapi.Message) {
	// Логика обработки команды /login
	reply := "Вы уже вошли в систему."
	b.sendMessage(msg.Chat.ID, reply)
}

func (b *AuthBot) handleUnknownCommand(msg *tgbotapi.Message) {
	reply := "Неизвестная команда. Доступные команды: /signup, /login"
	b.sendMessage(msg.Chat.ID, reply)
}

func (b *AuthBot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}
