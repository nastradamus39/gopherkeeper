package db

import (
	"log"

	"gophkeeper/gopherkeeper"

	"github.com/jmoiron/sqlx"
)

type repositories struct {
	Users   *UsersRepository
	Secrets *SecretsRepository
}

var repos *repositories

// InitDB инициализирует подключение к бд
func InitDB() (err error) {
	if gopherkeeper.DB, err = sqlx.Open("postgres", gopherkeeper.Cfg.DatabaseDsn); err != nil {
		log.Println(err)
	}

	repos = &repositories{
		Users: &UsersRepository{repo{
			table: "users",
			db:    gopherkeeper.DB,
		}},
		Secrets: &SecretsRepository{repo{
			table: "secrets",
			db:    gopherkeeper.DB,
		}},
	}

	return
}

// Repositories Возвращает список всех доступных репозиториев
func Repositories() *repositories {
	return repos
}
