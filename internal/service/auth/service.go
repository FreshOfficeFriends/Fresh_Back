package auth

import (
	"github.com/FreshOfficeFriends/SSO/internal/config"
	"github.com/FreshOfficeFriends/SSO/internal/domain"

	"time"
)

type Hasher interface {
	Hash(string) (string, error)
}

type UserRepo interface {
	UniqueEmail(email string) error
	SignUp(userInfo *domain.SignUp) error
	GetByCredentials(userInfo *domain.SignIn) (int, error)
	CreateRefreshToken(session domain.RefreshSession) error
	CredentialsByRefresh(refreshToken string) (int, time.Time, error)
}

type CacheUsers interface {
	SaveUser(hashEmail string, userInfo *domain.SignUp) error
	UserByHash(hashEmail string) ([]string, error)
	Exists(hashEmail string) bool
}

type Auth struct {
	hash       Hasher
	usersRepo  UserRepo
	cacheUsers CacheUsers
	jwt        *config.JWTConfig
}

func NewAuth(repo UserRepo, hash Hasher, cache CacheUsers, jwt *config.JWTConfig) *Auth {
	return &Auth{usersRepo: repo, hash: hash, cacheUsers: cache, jwt: jwt}
}
