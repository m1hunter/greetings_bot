package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username      string `gorm:"unique_index"`
	Password      string
	FirstName     string
	LastName      string
	Birthday      string         // Предполагаем, что день рождения хранится в формате "YYYY-MM-DD"
	Subscriptions []Subscription `gorm:"foreignkey:UserID"`
}

type Subscription struct {
	gorm.Model
	UserID       uint
	SubscribedTo uint
	Notification bool
}
