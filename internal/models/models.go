package models

import "time"

type User struct {
	ID        int
	Username  string
	Password  string
	FirstName string
	LastName  string
	Birthday  time.Time //хранить в date
	ChatID    int64     //bigint
}

type Subscription struct {
	ID                 int //id of
	UserID             int
	SubscribedTo       int
	isSendNotification bool
}
