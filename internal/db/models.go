package db

type User struct {
	Persist  bool   `db:"-" json:"-"`
	ID       int    `db:"id"`
	Login    string `db:"login" json:"login"`
	Password string `db:"password" json:"password"`
}
