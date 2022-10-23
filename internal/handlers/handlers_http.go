package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"gophkeeper/gopherkeeper"
	"gophkeeper/internal/db"

	"golang.org/x/crypto/bcrypt"
)

// RegisterHTTPHandler обработчик для регистрации пользователя
func RegisterHTTPHandler(w http.ResponseWriter, r *http.Request) {
	var user db.User

	// Обрабатываем входящий json
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Сохраняем пользователя в базе
	err := RegisterHandler(user)
	if err != nil {
		// Логин занят
		if errors.Is(err, gopherkeeper.ErrUserLoginConflict) {
			http.Error(w, gopherkeeper.ErrUserLoginConflict.Error(), http.StatusConflict)
		} else {
			InternalErrorHTTPResponse(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("пользователь успешно зарегистрирован"))
}

// LoginHTTPHandler авторизация пользователя
func LoginHTTPHandler(w http.ResponseWriter, r *http.Request) {
	var user db.User
	var err error

	// Обрабатываем входящий json
	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ищем пользователя в базе
	u, err := db.Repositories().Users.Find(user.Login)
	if err != nil {
		http.Error(w, "неверный логин/пароль", http.StatusUnauthorized)
		return
	}

	// проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "неверный логин/пароль", http.StatusUnauthorized)
		return
	}

	// Аунтефицируем пользователя
	err = AuthenticateUser(u, r, w)
	if err != nil {
		InternalErrorHTTPResponse(w, r, err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("authenticated"))
}

// InternalErrorHTTPResponse - возвращает пользователю 500 ошибку
func InternalErrorHTTPResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal server error. %s", err)
	http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
}

// AuthenticateUser создает сессию пользователя
func AuthenticateUser(user *db.User, r *http.Request, w http.ResponseWriter) error {
	// авторизуем пользователя
	session, err := gopherkeeper.SessionStore.Get(r, gopherkeeper.SessionName)
	if err != nil {
		return err
	}
	session.Values["userId"] = user.Login
	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}
