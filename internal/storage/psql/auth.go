package psql

import (
	"database/sql"

	"github.com/FreshOfficeFriends/SSO/internal/domain"
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{db: db}
}

func (a *Users) UniqueEmail(email string) error {
	return a.db.QueryRow("SELECT email from users where email=$1", email).Scan()
}

func (a *Users) SignUp(user *domain.SignUp) error {
	_, err := a.db.Exec("INSERT INTO users (first_name, second_name, email, password, birthday) values ($1, $2, $3, $4, $5)",
		user.FirstName, user.SecondName, user.Email, user.Password, user.Birthday)

	return err
}
