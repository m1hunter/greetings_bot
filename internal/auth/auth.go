package auth

import (
	"projectik/internal/database"
	"projectik/internal/models"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(username, password, firstName, lastName, birthday string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return errors.Wrap(err, "failed to hash password")
	}

	user := models.User{
		Username:  username,
		Password:  hashedPassword,
		FirstName: firstName,
		LastName:  lastName,
		Birthday:  birthday,
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to create user")
	}

	return nil
}

func AuthenticateUser(username, password string) (*models.User, bool) {
	var user models.User
	result := database.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, false
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, false
	}

	return &user, true
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
