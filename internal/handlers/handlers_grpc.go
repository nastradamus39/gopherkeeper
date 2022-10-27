package handlers

import (
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gophkeeper/internal/db"
)

// LoginGRPCHandler авторизация пользователя
func LoginGRPCHandler(login string, pass string) (token string, err error) {
	// Ищем пользователя в базе
	u, err := db.Repositories().Users.Find(login)
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, "Неверный логин/пароль")
	}

	// проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass))
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, "Неверный логин/пароль")
	}

	return u.Token, nil
}
