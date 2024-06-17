package telegram

import (
	"fmt"
	"log"
	"strings"

	"projectik/internal/models"
	"projectik/internal/notification"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type Bot struct {
	api          *tgbotapi.BotAPI
	notification *notification.NotificationService
	db           *gorm.DB
}

func NewBot(token string, db *gorm.DB) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:          bot,
		notification: notification.NewNotificationService(db),
		db:           db,
	}, nil
}

func (b *Bot) Start() error {
	log.Println("Telegram bot started")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		go b.handleUpdate(update)
	}

	return nil
}

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	if update.Message.IsCommand() {
		b.handleCommand(update.Message)
	} else {
		b.handleTextMessage(update.Message)
	}
}

func (b *Bot) handleCommand(msg *tgbotapi.Message) {
	switch msg.Command() {
	case "start":
		b.handleStart(msg)
	case "signup":
		b.handleSignup(msg)
	case "login":
		b.handleLogin(msg)
	case "list":
		b.handleList(msg)
	case "sublist":
		b.handleSubscribedList(msg)
	case "sub":
		b.handleSubscribe(msg)
	case "unsub":
		b.handleUnsubscribe(msg)
	case "showcommands":
		b.handleShowCommands(msg)
	default:
		b.handleUnknownCommand(msg)
	}
}

func (b *Bot) handleStart(msg *tgbotapi.Message) {
	reply := "Привет! Я бот для уведомлений о днях рождения. Используй /signup для регистрации."
	b.sendMessage(msg.Chat.ID, reply)
}

func (b *Bot) handleSignup(msg *tgbotapi.Message) {
	args := msg.CommandArguments()
	if args == "" {
		reply := "Для регистрации введи /signup <username> <password> <firstname> <lastname> <birthday (YYYY-MM-DD)>"
		b.sendMessage(msg.Chat.ID, reply)
		return
	}

	params := strings.Fields(args)
	if len(params) < 5 {
		reply := "Неверный формат. Используй /signup <username> <password> <firstname> <lastname> <birthday (YYYY-MM-DD)>"
		b.sendMessage(msg.Chat.ID, reply)
		return
	}

	username := params[0]
	password := params[1]
	firstname := params[2]
	lastname := params[3]
	birthday := params[4]

	user := models.User{
		Username:  username,
		Password:  password,
		FirstName: firstname,
		LastName:  lastname,
		Birthday:  birthday,
	}

	err := b.db.Create(&user).Error
	if err != nil {
		reply := "Ошибка при регистрации пользователя: " + err.Error()
		b.sendMessage(msg.Chat.ID, reply)
		return
	}

	reply := "Пользователь зарегистрирован успешно!"
	b.sendMessage(msg.Chat.ID, reply)
}

func (b *Bot) handleLogin(msg *tgbotapi.Message) {
	reply := "Функция входа пока не реализована."
	b.sendMessage(msg.Chat.ID, reply)
}

func (b *Bot) handleList(msg *tgbotapi.Message) {
	var users []models.User
	err := b.db.Find(&users).Error
	if err != nil {
		reply := "Ошибка при получении списка сотрудников: " + err.Error()
		b.sendMessage(msg.Chat.ID, reply)
		return
	}

	reply := "Список сотрудников:\n"
	for _, user := range users {
		reply += fmt.Sprintf("ID: %d, Имя: %s %s\n", user.ID, user.FirstName, user.LastName)
	}
	b.sendMessage(msg.Chat.ID, reply)
}

func (b *Bot) handleSubscribedList(msg *tgbotapi.Message) {
	var user models.User
	if err := b.db.Where("username = ?", msg.From.UserName).First(&user).Error; err != nil {
		b.sendMessage(msg.Chat.ID, "Ошибка при получении данных пользователя.")
		return
	}

	var subscriptions []models.Subscription
	if err := b.db.Where("user_id = ?", user.ID).Find(&subscriptions).Error; err != nil {
		b.sendMessage(msg.Chat.ID, "Ошибка при получении подписок.")
		return
	}

	reply := "Список подписок:\n"
	for _, sub := range subscriptions {
		var subscribedTo models.User
		if err := b.db.First(&subscribedTo, sub.SubscribedTo).Error; err == nil {
			reply += fmt.Sprintf("ID: %d, Имя: %s %s\n", subscribedTo.ID, subscribedTo.FirstName, subscribedTo.LastName)
		}
	}

	b.sendMessage(msg.Chat.ID, reply)
}

func (b *Bot) handleSubscribe(msg *tgbotapi.Message) {
	args := msg.CommandArguments()
	if args == "" {
		b.sendMessage(msg.Chat.ID, "Используй /sub <user_id>")
		return
	}

	var user models.User
	if err := b.db.Where("username = ?", msg.From.UserName).First(&user).Error; err != nil {
		b.sendMessage(msg.Chat.ID, "Ошибка при получении данных пользователя.")
		return
	}

	subscribedToID := 0
	fmt.Sscanf(args, "%d", &subscribedToID)

	var subscribedTo models.User
	if err := b.db.First(&subscribedTo, subscribedToID).Error; err != nil {
		b.sendMessage(msg.Chat.ID, "Пользователь с указанным ID не найден.")
		return
	}

	subscription := models.Subscription{
		UserID:       user.ID,
		SubscribedTo: subscribedTo.ID,
	}

	if err := b.db.Create(&subscription).Error; err != nil {
		b.sendMessage(msg.Chat.ID, "Ошибка при подписке.")
		return
	}

	b.sendMessage(msg.Chat.ID, "Подписка успешно оформлена!")
}

func (b *Bot) handleUnsubscribe(msg *tgbotapi.Message) {
	args := msg.CommandArguments()
	if args == "" {
		b.sendMessage(msg.Chat.ID, "Используй /unsub <user_id>")
		return
	}

	var user models.User
	if err := b.db.Where("username = ?", msg.From.UserName).First(&user).Error; err != nil {
		b.sendMessage(msg.Chat.ID, "Ошибка при получении данных пользователя.")
		return
	}

	subscribedToID := 0
	fmt.Sscanf(args, "%d", &subscribedToID)

	var subscription models.Subscription
	if err := b.db.Where("user_id = ? AND subscribed_to = ?", user.ID, subscribedToID).First(&subscription).Error; err != nil {
		b.sendMessage(msg.Chat.ID, "Подписка не найдена.")
		return
	}

	if err := b.db.Delete(&subscription).Error; err != nil {
		b.sendMessage(msg.Chat.ID, "Ошибка при отписке.")
		return
	}

	b.sendMessage(msg.Chat.ID, "Отписка успешно оформлена!")
}

func (b *Bot) handleUnknownCommand(msg *tgbotapi.Message) {
	reply := "Неизвестная команда. Используй /start для начала работы."
	b.sendMessage(msg.Chat.ID, reply)
}

func (b *Bot) handleTextMessage(msg *tgbotapi.Message) {
	reply := "Я понимаю только команды. Используй /start для начала работы."
	b.sendMessage(msg.Chat.ID, reply)
}

func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func (b *Bot) handleShowCommands(msg *tgbotapi.Message) {
	reply := "Доступные команды:\n" +
		"/start - начать работу с ботом\n" +
		"/signup <username> <password> <firstname> <lastname> <birthday (YYYY-MM-DD)> - зарегистрироваться\n" +
		"/list - показать список пользователей\n" +
		"/sublist - показать список подписок\n" +
		"/sub <username> - подписаться на пользователя\n" +
		"/unsub <username - отписаться от пользователя\n"
	b.sendMessage(msg.Chat.ID, reply)
}
