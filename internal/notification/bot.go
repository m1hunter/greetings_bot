package notification

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"projectik/internal/auth"
	"projectik/internal/models"
)

type Bot struct {
	api          *tgbotapi.BotAPI
	notification *NotificationService
	auth         *auth.AuthBot
	db           *sql.DB
}

// чтобы отправить сообщение определенному пользователю, будем использовать Message.Chat.ID,
// Message - огромная структура, в которой есть все для работы с сообщениями и его
// отправителем
func NewBot(token string, db *sql.DB) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	log.Printf("Authorized on account %s", botAPI.Self.UserName)

	bot := &Bot{
		api: botAPI,
		db:  db,
	}

	bot.auth = auth.NewAuthBot(botAPI, db)

	bot.notification = NewNotificationService(db, bot)

	bot.notification.StartBirthdayNotification()

	return bot, nil
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
	case "signup", "whoami": // Передача команд авторизации в AuthBot
		b.auth.HandleMessage(tgbotapi.Update{Message: msg})
	case "list":
		b.handleShowList(msg)
	case "sublist":
		b.handleShowSubscribedList(msg)
	case "sub":
		b.handleSubscribe(msg)
	case "unsub":
		b.handleUnsubscribe(msg)
	case "help":
		b.handleHelp(msg)
	default:
		b.handleUnknownCommand(msg)
	}
}

func (b *Bot) handleStart(msg *tgbotapi.Message) {
	query := `UPDATE users SET chat_id = $1 WHERE username = $2`
	_, err := b.db.Exec(query, msg.Chat.ID, msg.From.UserName)
	if err != nil {
		log.Printf("Error updating chat_id for user %s: %v", msg.From.UserName, err)
	}

	reply :=
		"Привет! Я бот для уведомлений о днях рождениях. " +
			"Используй /signup для регистрации и /whoami " +
			"для получения информации о вашем аккаунте.\n" +
			"Для получения списка команд введите /help."
	b.SendMessage(msg.Chat.ID, reply)
}

func (b *Bot) handleShowList(msg *tgbotapi.Message) {
	query := "SELECT id, username, firstname, lastname FROM users"
	rows, err := b.db.Query(query)
	if err != nil {
		b.SendMessage(msg.Chat.ID, "При получении списка пользователей произошла ошибка, повторите попытку позже.")
		log.Printf("Error fetching user list: %v", err)
		return
	}
	defer rows.Close()

	reply := "Список пользователей:\n"
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName); err != nil {
			b.SendMessage(msg.Chat.ID, "При получении данных пользователей произошла ошибка, повторите попытку позже.")
			log.Printf("Error scanning user data: %v", err)
			return
		}
		reply += fmt.Sprintf("ID: %d, Имя: %s %s\n", user.ID, user.FirstName, user.LastName)
	}
	b.SendMessage(msg.Chat.ID, reply)
	log.Printf("User list sent to chat %d", msg.Chat.ID)
}

func (b *Bot) handleShowSubscribedList(msg *tgbotapi.Message) {
	var user models.User
	if err := b.db.QueryRow("SELECT id FROM users WHERE username = $1", msg.From.UserName).Scan(&user.ID); err != nil {
		b.SendMessage(msg.Chat.ID, "При получении данных пользователя произошла ошибка, повторите попытку позже.")
		log.Printf("Error fetching user data for %s: %v", msg.From.UserName, err)
		return
	}

	query := "SELECT subscribed_to FROM subscriptions WHERE user_id = $1"
	rows, err := b.db.Query(query, user.ID)
	if err != nil {
		b.SendMessage(msg.Chat.ID, "При получении списка подписок произошла ошибка, повторите попытку позже.")
		log.Printf("Error fetching subscriptions for user %d: %v", user.ID, err)
		return
	}
	defer rows.Close()

	reply := "Список подписок:\n"
	for rows.Next() {
		var subscribedToID int
		if err := rows.Scan(&subscribedToID); err != nil {
			b.SendMessage(msg.Chat.ID, "При сканировании Ваших подписок произошла ошибка, повторите попытку позже.")
			log.Printf("Error scanning subscription for user %d: %v", user.ID, err)
			return
		}

		var subscribedTo models.User
		if err := b.db.QueryRow("SELECT id, firstname, lastname FROM users WHERE id = $1", subscribedToID).Scan(&subscribedTo.ID, &subscribedTo.FirstName, &subscribedTo.LastName); err == nil {
			reply += fmt.Sprintf("ID: %d, Имя: %s %s\n", subscribedTo.ID, subscribedTo.FirstName, subscribedTo.LastName)
		}
	}

	b.SendMessage(msg.Chat.ID, reply)
	log.Printf("Subscription list sent to user %s", msg.Chat.ID)
	// Здесь можно было использовать и msg.From.UserName, но в логах консоли лучше видеть ID
	// чата, а не никнейм пользователя
}

func (b *Bot) handleSubscribe(msg *tgbotapi.Message) {
	args := msg.CommandArguments()
	if args == "" {
		b.SendMessage(msg.Chat.ID, "Используйте /sub <user_id>")
		return
	}

	var user models.User
	if err := b.db.QueryRow("SELECT id FROM users WHERE username = $1", msg.From.UserName).Scan(&user.ID); err != nil {
		b.SendMessage(msg.Chat.ID, "При получении данных пользователя произошла ошибка, повторите попытку позже.")
		log.Printf("Error fetching data for user %s: %v", msg.From.UserName, err)
		return
	}

	subscribedToID := 0
	fmt.Sscanf(args, "%d", &subscribedToID)

	var subscribedTo models.User
	if err := b.db.QueryRow("SELECT id FROM users WHERE id = $1", subscribedToID).Scan(&subscribedTo.ID); err != nil {
		b.SendMessage(msg.Chat.ID, "Пользователь с указанным ID не найден.")
		log.Printf("User with ID %d not found for subscription", subscribedToID)
		return
	}

	query := `INSERT INTO subscriptions (user_id, subscribed_to, is_send_notification) 
	          VALUES ($1, $2, false)`
	_, err := b.db.Exec(query, user.ID, subscribedTo.ID)
	if err != nil {
		b.SendMessage(msg.Chat.ID, "При попытке подписки на уведомления "+
			"произошла ошибка, повторите попытку позже.")
		log.Printf("Error during subscription: %v", err)
		return
	}

	b.SendMessage(msg.Chat.ID, "Подписка успешно оформлена!")
	log.Printf("User %d subscribed to user %d", user.ID, subscribedTo.ID)
}

func (b *Bot) handleUnsubscribe(msg *tgbotapi.Message) {
	args := msg.CommandArguments()
	if args == "" {
		b.SendMessage(msg.Chat.ID, "Используйте /unsub <user_id>")
		return
	}

	var user models.User
	if err := b.db.QueryRow("SELECT id FROM users WHERE username = $1", msg.From.UserName).Scan(&user.ID); err != nil {
		b.SendMessage(msg.Chat.ID, "При получении данных пользователя произошла ошибка, повторите попытку позже.")
		log.Printf("Error fetching data for user %s: %v", msg.From.UserName, err)
		return
	}

	subscribedToID := 0
	fmt.Sscanf(args, "%d", &subscribedToID)

	query := "DELETE FROM subscriptions WHERE user_id = $1 AND subscribed_to = $2"
	_, err := b.db.Exec(query, user.ID, subscribedToID)
	if err != nil {
		b.SendMessage(msg.Chat.ID,
			"При отписке от уведомлений на пользователя произошла ошибка, "+
				"повторите попытку позже.")
		log.Printf("Error during unsubscription: %v", err)
		return
	}

	b.SendMessage(msg.Chat.ID, "Отписка успешно оформлена!")
	log.Printf("User %d unsubscribed from user %d", user.ChatID, subscribedToID)
}

func (b *Bot) handleUnknownCommand(msg *tgbotapi.Message) {
	reply := "Неизвестная команда. Используйте /start для начала работы или /help для получения списка команд."
	b.SendMessage(msg.Chat.ID, reply)
}

func (b *Bot) handleTextMessage(msg *tgbotapi.Message) {
	reply := "Я понимаю только команды. Используйте /start для начала работы или /help для получения списка команд."
	b.SendMessage(msg.Chat.ID, reply)
}

func (b *Bot) handleHelp(msg *tgbotapi.Message) {
	reply :=
		"Список команд:\n" +
			"/signup - регистрация пользователя\n" +
			"/whoami - информация о текущем логине\n" +
			"/sub - подписаться на уведомления о дне рождении пользователя\n" +
			"/unsub - отписаться от уведомления о дне рождении пользователя\n" +
			"/list - список всех пользователей\n" +
			"/sublist - список пользователей, на которых Вы подписаны"
	b.SendMessage(msg.Chat.ID, reply)
}

func (b *Bot) SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}
