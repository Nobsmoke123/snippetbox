package mocks

import (
	"time"

	"github.com/Nobsmoke123/snippetbox/internal/models"
)

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int,string, error) {
	if email == "alice@gmail.com" && password == "pa$$word" {
		return 1, "admin",nil
	}
	return 0,"",models.ErrInvalidCredentials
}

func (m *UserModel) Get(id int) (models.User, error){
	if id == 1 {
		user := models.User{
			ID: 1,
			Email: "alice@gmail.com",
			Name: "Alice",
			Created: time.Now(),
		}

		return user, nil
	}

	return models.User{}, models.ErrNoRecord
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) PasswordUpdate(id int, currentPassword, newPassword string) error {
	if id == 1 {
		return nil
	}else{
		return models.ErrNoRecord
	}
}