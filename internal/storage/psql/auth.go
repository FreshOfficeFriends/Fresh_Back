package psql

import (
	"database/sql"
	"time"

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

func (a *Users) GetByCredentials(user *domain.SignIn) (int, error) {
	var id int
	return id, a.db.QueryRow("SELECT id from users where email=$1 and password=$2", user.Email, user.Password).Scan(&id)
}

func (a *Users) CreateRefreshToken(inp domain.RefreshSession) error {
	_, err := a.db.Exec("INSERT INTO refresh_tokens (user_id, token, expires_at) values ($1, $2, $3)", inp.UserID,
		inp.Token, inp.ExpiresAt)

	return err
}

func (a *Users) CredentialsByRefresh(refreshToken string) (int, time.Time, error) {
	var id int
	var exp time.Time
	return id, exp, a.db.QueryRow("SELECT user_id, expires_at from refresh_tokens where token=$1", refreshToken).Scan(&id, &exp)
}
