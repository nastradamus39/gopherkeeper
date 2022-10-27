package db

import (
	"errors"
	"fmt"

	"gophkeeper/gopherkeeper"

	"github.com/jmoiron/sqlx"
)

type repository interface {
	// Save сохраняет сущность в бд
	Save(entity interface{}) error
	// Delete удаляет сущность из бд
	Delete(entity interface{}) error
}

type repo struct {
	table string
	db    *sqlx.DB
}

type UsersRepository struct {
	repo
}

// Save сохраняет пользователя
func (r *UsersRepository) Save(user interface{}) error {
	u, ok := user.(*User)
	if !ok {
		return errors.New("unsupported type")
	}

	if !u.Persist {
		res, err := r.db.NamedQuery(`INSERT INTO users(login, password, token) 
			VALUES (:login, :password, :token) on conflict (login) DO NOTHING RETURNING login`, &u)

		if err != nil && res.Err() != nil {
			return err
		}

		if !res.Next() {
			return fmt.Errorf("%w", gopherkeeper.ErrUserLoginConflict)
		}
	}

	return nil
}

// Delete удаляет пользователя в базе
func (r *UsersRepository) Delete(user interface{}) error {
	return nil
}

// Find поиск пользователя по логину
func (r *UsersRepository) Find(login string) (user *User, err error) {
	user = &User{}
	err = r.db.Get(user, "SELECT * FROM users WHERE login = $1", login)
	user.Persist = true
	return
}

// FindByToken ищет пользователя по токену
func (r *UsersRepository) FindByToken(token string) (user *User, err error) {
	user = &User{}
	err = r.db.Get(user, "SELECT * FROM users WHERE token = $1", token)
	user.Persist = true
	return
}
