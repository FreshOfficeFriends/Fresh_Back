package auth

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/FreshOfficeFriends/SSO/internal/domain"
)

func (a *Auth) UniqueEmail(email string) error {
	if err := a.usersRepo.UniqueEmail(email); err == sql.ErrNoRows {
		return nil
	}
	return errors.New("email not unique")
}

func (a *Auth) SignUp(userInfo *domain.SignUp) error {
	hashEmail, err := a.hash.Hash(userInfo.Email)
	if err != nil {
		return err
	}

	if err = a.cacheUsers.SaveUser(hashEmail, userInfo); err != nil {
		return err
	}

	return SendEmail(userInfo.Email, hashEmail)
}

func (a *Auth) SignIn(userInfo *domain.SignIn) (string, string, error) {
	hashPass, err := a.hash.Hash(userInfo.Password)
	if err != nil {
		return "", "", err
	}

	userInfo.Password = hashPass

	id, err := a.usersRepo.GetByCredentials(userInfo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", domain.ErrUserNotFound
		}
		return "", "", err
	}

	return a.generateTokens(id)
}

func (a *Auth) SaveUser(hashEmail string) error {
	if !a.cacheUsers.Exists(hashEmail) {
		return errors.New("bye bye")
	}

	data, err := a.cacheUsers.UserByHash(hashEmail)
	if err != nil || len(data) == 0 {
		return errors.New("redis err")
	}

	userInfo, err := a.parseUserInfo(data[0])
	if err != nil {
		return err
	}

	return a.usersRepo.SignUp(userInfo)
}

func (a *Auth) parseUserInfo(userInfoFromRedis string) (*domain.SignUp, error) {
	userInfo := new(domain.SignUp)
  
	user := strings.Split(userInfoFromRedis, " ")

	userInfo.FirstName = user[0]
	userInfo.SecondName = user[1]
	userInfo.Email = user[2]
	userInfo.Birthday = user[3]

	hashPass, err := a.hash.Hash(user[4])
	if err != nil {
		return nil, err
	}

	userInfo.Password = hashPass

	return userInfo, nil
}
