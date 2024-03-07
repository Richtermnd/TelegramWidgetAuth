package domain

const userTable = "users"

var userColumns = []string{
	"tg_id",
	"first_name",
	"last_name",
	"username",
	"photo_url",
	"last_login",
	"registered",
}

type User struct {
	TGID       int64  `db:"tg_id" json:"tg_id"`
	FirstName  string `db:"first_name" json:"first_name"`
	LastName   string `db:"last_name" json:"last_name"`
	Username   string `db:"username" json:"username"`
	PhotoURL   string `db:"photo_url" json:"photo_url"`
	LastLogin  int64  `db:"last_login" json:"last_login"` // UNIX timestamp
	Registered int64  `db:"registered" json:"registered"` // UNIX timestamp
}

func (u *User) Table() string {
	return userTable
}

func (u *User) Columns() []string {
	return userColumns[:]
}

func (u *User) Values() []interface{} {
	return []interface{}{
		u.TGID,
		u.FirstName,
		u.LastName,
		u.Username,
		u.PhotoURL,
		u.LastLogin,
		u.Registered,
	}
}

func (u *User) Map() map[string]interface{} {
	return map[string]interface{}{
		"tg_id":      u.TGID,
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"username":   u.Username,
		"photo_url":  u.PhotoURL,
		"last_login": u.LastLogin,
		"Registered": u.Registered,
	}
}
