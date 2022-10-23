package main

import (
	"embed"
	"log"
	"net/http"
	"os"

	"gophkeeper/gopherkeeper"
	"gophkeeper/internal/db"
	"gophkeeper/internal/handlers"

	_ "github.com/lib/pq"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	r := Router()

	// Logger
	flog, err := os.OpenFile(`server.log`, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer flog.Close()

	//log.SetOutput(flog)

	// Переменные окружения в конфиг
	err = env.Parse(&gopherkeeper.Cfg)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Инициализация подключения к бд
	err = db.InitDB()
	if err != nil {
		log.Fatal(err)
		return
	}

	// миграции
	goose.SetBaseFS(embedMigrations)
	if err = goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
		return
	}
	if err = goose.Up(gopherkeeper.DB.DB, "migrations"); err != nil {
		log.Fatal(err)
		return
	}

	// запускаем сервер
	err = http.ListenAndServe("localhost:8080", r)
	if err != nil {
		log.Printf("Не удалось запустить сервер. %s", err)
		return
	}
}

func Router() (r *chi.Mux) {
	r = chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	r.Post("/user/register", handlers.RegisterHTTPHandler)
	r.Post("/user/login", handlers.LoginHTTPHandler)

	// закрытые авторизацией эндпоинты
	r.Mount("/user", privateRouter())

	return r
}

// privateRouter Роутер для закрытых авторизацией эндпоинтов
func privateRouter() http.Handler {
	r := chi.NewRouter()
	//r.Use(middlewares.UserAuth) // проверка авторизации

	return r
}
