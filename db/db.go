package db

import (
	"errors"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/rooty123/libs/dbhandler"
)

type UsersDBHandler struct {
	dbhandler.DBHandler
}

type User struct {
	ID           int       `pg:"id,pk"`
	ChatID       int64     `pg:"chat_id,unique,notnull"`
	TelegramID   int64     `pg:"telegram_id,unique,notnull"`
	FirstName    string    `pg:"first_name"`
	LastName     string    `pg:"last_name"`
	LanguageCode string    `pg:"language_code"`
	Username     string    `pg:"username"`
	State        string    `pg:"state"`
	CreationTime time.Time `pg:"creation_time,default:now()"`
}

func (udb *UsersDBHandler) CreateUser(chatID, telegramID int64, firstName, lastName, languageCode, username string) error {
	user := &User{
		ChatID:       chatID,
		TelegramID:   telegramID,
		FirstName:    firstName,
		LastName:     lastName,
		LanguageCode: languageCode,
		Username:     username,
		CreationTime: time.Now(),
	}
	_, err := udb.Conn.Model(user).Insert()
	return err
}

func (udb *UsersDBHandler) GetUser(chatID int64) (*User, error) {
	user := new(User)
	err := udb.Conn.Model(user).Where("chat_id = ?", chatID).Select()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil // Return nil if user not found
		}
		return nil, err
	}
	return user, nil
}

func (udb *UsersDBHandler) UpdateUser(chatID int64, firstName, lastName, languageCode, username string) error {
	_, err := udb.Conn.Model(&User{
		FirstName:    firstName,
		LastName:     lastName,
		LanguageCode: languageCode,
		Username:     username,
	}).Where("chat_id = ?", chatID).Update()
	return err
}

func (udb *UsersDBHandler) ListUsers() ([]User, error) {
	var users []User
	err := udb.Conn.Model(&users).Select()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (udb *UsersDBHandler) DeleteUser(chatID int64) error {
	_, err := udb.Conn.Model(&User{}).Where("chat_id = ?", chatID).Delete()
	return err
}
