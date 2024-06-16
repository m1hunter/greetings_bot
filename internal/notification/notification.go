package notification

import (
	"fmt"
	"log"
	"projectik/internal/models"
	"time"

	"gorm.io/gorm"
)

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

func (n *NotificationService) CheckBirthdays() {
	var users []models.User
	if err := n.db.Find(&users).Error; err != nil {
		log.Printf("Error fetching users: %v", err)
		return
	}

	today := time.Now().Format("2006-01-02")

	for _, user := range users {
		if user.Birthday == today {
			n.notifySubscribers(user)
		}
	}
}

func (n *NotificationService) notifySubscribers(user models.User) {
	var subscriptions []models.Subscription
	err := n.db.Where("subscribed_to = ?", user.ID).Find(&subscriptions).Error
	if err != nil {
		log.Printf("Error fetching subscriptions: %v", err)
		return
	}

	for _, subscription := range subscriptions {
		var subscriber models.User
		if err := n.db.Where("id = ?", subscription.UserID).First(&subscriber).Error; err == nil {
			message := fmt.Sprintf("Сегодня день рождения у %s %s!", user.FirstName, user.LastName)
			n.sendMessage(subscriber, message)
		}
	}
}

func (n *NotificationService) sendMessage(user models.User, message string) {
	log.Printf("Sending message to %s: %s", user.Username, message)
}
