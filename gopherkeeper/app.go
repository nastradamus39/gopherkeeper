package gopherkeeper

import (
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

const (
	SessionName               = "gopherKeeperSid"
	TokenPath                 = ".gph_key"
	ContextUserKey ContextKey = iota
)

type ContextKey int8

// DB подключение к базе
var DB *sqlx.DB

// SessionStore хранилище сессий
var SessionStore = sessions.NewCookieStore([]byte("aaa"))

// Cfg конфиг приложения
var Cfg Config

// Config конфиг приложения
type Config struct {
	ServerAddress string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	DatabaseDsn   string `env:"DATABASE_URI" envDefault:""`
}
