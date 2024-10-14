package notification

import (
	"database/sql"
	"fmt"
	"log"
	"projectik/internal/models" // Импортируем модели
	"time"
)

type NotificationService struct {
	db  *sql.DB
	bot *Bot
}

func NewNotificationService(db *sql.DB, bot *Bot) *NotificationService {
	return &NotificationService{
		db:  db,
		bot: bot,
	}
}

// Метод StartBirthdayNotification() для запуска фоновой горутины проверки дней рождения,
// метод CheckBirthdays() срабатывает в 00:00 по серверу по МСК, но это сделано через мини-костыль
func (n *NotificationService) StartBirthdayNotification() {
	go func() {
		for {
			now := time.Now()
			// Находим время до следующего 00:00 по MSK.
			//TODO refactor to cron
			nextMidnight := now.Truncate(21 * time.Hour).Add(24 * time.Hour)
			timeUntilNextMidnight := nextMidnight.Sub(now)

			log.Printf("The next check for birthdays will be %s", nextMidnight)
			time.Sleep(timeUntilNextMidnight)

			n.CheckBirthdays()
		}
	}()
}

func (n *NotificationService) CheckBirthdays() {
	query := "SELECT id, username, firstname, lastname, birthday FROM users"
	rows, err := n.db.Query(query)
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		return
	}
	defer rows.Close()

	today := time.Now()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Birthday); err != nil {
			log.Printf("Error scanning user: %v", err)
			continue
		}

		// Уведомление пользователей, если совпал месяц и день
		if user.Birthday.Month() == today.Month() && user.Birthday.Day() == today.Day() {
			n.notifySubscribers(user)
		}
	}
}

func (n *NotificationService) notifySubscribers(user models.User) {
	query := "SELECT user_id FROM subscriptions WHERE subscribed_to = $1 AND is_send_notification = false"
	rows, err := n.db.Query(query, user.ID)
	if err != nil {
		log.Printf("Error getting subscriptions: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var subscriberID int
		if err := rows.Scan(&subscriberID); err != nil {
			log.Printf("Error scanning subscription: %v", err)
			continue
		}

		// Получаем chat_id подписчика
		var subscriber models.User
		err := n.db.QueryRow("SELECT chat_id FROM users WHERE id = $1", subscriberID).Scan(&subscriber.ChatID)
		if err != nil {
			log.Printf("Error getting chat_id for subscriber, error: %v", err)
			continue
		}

		// Отправляем уведомление подписчику о дне рождения пользователя
		message := fmt.Sprintf("Сегодня день рождения у %s %s!", user.FirstName, user.LastName)
		n.bot.SendMessage(subscriber.ChatID, message) // Используем chat_id для отправки сообщения

		// Обновляем запись, чтобы больше не отправлять уведомления
		_, err = n.db.Exec("UPDATE subscriptions SET is_send_notification = true WHERE user_id = $1 AND subscribed_to = $2", subscriberID, user.ID)
		if err != nil {
			log.Printf("Error updating subscription notification, error: %v", err)
		}
	}
}
