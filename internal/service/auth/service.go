package auth

import "github.com/FreshOfficeFriends/SSO/internal/domain"

type Hasher interface {
	Hash(string) (string, error)
}

type UserRepo interface {
	UniqueEmail(email string) error
	SignUp(userInfo *domain.SignUp) error
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
}

func NewAuth(repo UserRepo, hash Hasher, cache CacheUsers) *Auth {
	return &Auth{usersRepo: repo, hash: hash, cacheUsers: cache}
}
