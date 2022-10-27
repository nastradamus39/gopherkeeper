package db

type User struct {
	Persist  bool   `db:"-" json:"-"`
	ID       int    `db:"id"`
	Login    string `db:"login" json:"login"`
	Password string `db:"password" json:"password"`
	Token    string `db:"token" json:"-"`
}

type Secret struct {
	Persist bool   `db:"-" json:"-"`
	ID      int    `db:"id"`
	Login   string `db:"login" json:"login"`
	Comment string `db:"comment" json:"comment"`
	Card    string `db:"card" json:"-"`
	Attach  string `db:"attach" json:"-"`
}
