package handlers

import (
	"gophkeeper/internal/db"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(user db.User) error {
	// Считаем хеш пароля для дальнейшего сохранения
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	// Сохраняем пользователя в базу
	err = db.Repositories().Users.Save(&user)
	if err != nil {
		return err
	}

	return nil
}
