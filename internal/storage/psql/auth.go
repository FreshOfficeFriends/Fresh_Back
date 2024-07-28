package psql

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/FreshOfficeFriends/SSO/pkg/logger"

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
	uuid := uuid.New()

	_, err := a.db.Exec("INSERT INTO users (uuid, first_name, second_name, email, password, birthday) values ($1, $2, $3, $4, $5, $6)",
		uuid, user.FirstName, user.SecondName, user.Email, user.Password, user.Birthday)

	return err
}

func (a *Users) GetByCredentials(user *domain.SignIn) (int, error) {
	var id int
	return id, a.db.QueryRow("SELECT id from users where email=$1 and password=$2", user.Email, user.Password).Scan(&id)
}

func (a *Users) CreateRefreshToken(inp domain.RefreshSession) error {
	//немного БЗ в репо слое)(
	_, err := a.db.Exec("DELETE FROM refresh_tokens where user_id=$1", inp.UserID)
	if err != nil {
		return err
	}
	_, err = a.db.Exec("INSERT INTO refresh_tokens (user_id, token, expires_at) values ($1, $2, $3)", inp.UserID,
		inp.Token, inp.ExpiresAt)

	return err
}

func (a *Users) CredentialsByRefresh(refreshToken string) (int, time.Time, error) {
	var id int
	var exp time.Time
	return id, exp, a.db.QueryRow("SELECT user_id, expires_at from refresh_tokens where token=$1", refreshToken).Scan(&id, &exp)
}

func (a *Users) UUID(email string) (string, error) {
	var uuid string
	return uuid, a.db.QueryRow("SELECT uuid from users where email=$1", email).Scan(&uuid)
}

func (a *Users) ChangePass(uuid string, hashPass string) error {
	logger.Debug("change pass", zap.String("uuid", uuid), zap.String("hashPass", hashPass))
	_, err := a.db.Exec("UPDATE users SET password=$1 WHERE uuid=$2", hashPass, uuid)
	return err
}

//func (a *Users) ChangePass(uuid string, hashPass string) error {
//	_, err := a.db.Exec("UPDATE users set password='$1' where uuid='$2'", hashPass, uuid)
//	return err
//}
