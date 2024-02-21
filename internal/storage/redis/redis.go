package rds

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/FreshOfficeFriends/SSO/internal/domain"
)

type Client struct {
	db *redis.Client
}

func New(db *redis.Client) *Client {
	return &Client{db: db}
}

func (c *Client) SaveUser(hashEmail string, userInfo *domain.SignUp) error {
	user := fmt.Sprintf("%s %s %s %s %s", userInfo.FirstName, userInfo.SecondName, userInfo.Email,
		userInfo.Birthday, userInfo.Password)
	err := c.db.SAdd(context.Background(), hashEmail, user).Err()
	if err != nil {
		return err
	}
	return c.db.Expire(context.Background(), hashEmail, time.Hour).Err()
}

func (c *Client) UserByHash(hashEmail string) ([]string, error) {
	data := c.db.SMembers(context.Background(), hashEmail)
	return data.Val(), data.Err()
}

func (c *Client) Exists(hashEmail string) bool {
	if ok := c.db.Exists(context.Background(), hashEmail).Val(); ok == 0 {
		return false
	}
	return true
}
