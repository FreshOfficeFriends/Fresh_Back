package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"github.com/FreshOfficeFriends/SSO/internal/domain"
	"github.com/FreshOfficeFriends/SSO/pkg/logger"
)

func (a *Auth) generateTokens(userId int) (string, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"UserId":    strconv.Itoa(userId),
		"IssuedAt":  time.Now().Unix(),
		"ExpiresAt": time.Now().Add(a.jwt.AccessTTL).Unix(),
	})

	accessToken, err := token.SignedString(a.jwt.Secret)
	if err != nil {
		return "", "", err
	}

	logger.Debug("", zap.String("access token", accessToken))

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}

	logger.Debug("", zap.String("refresh token", refreshToken))

	if err = a.usersRepo.CreateRefreshToken(domain.RefreshSession{
		UserID:    userId,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, err

}

func newRefreshToken() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", bytes), err
}

func (a *Auth) ParseToken(token string) (int, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return a.jwt.Secret, nil
	})

	if err != nil {
		return 0, err
	}

	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	subject, ok := claims["UserId"].(string)
	if !ok {
		return 0, errors.New("invalid subject")
	}

	if claims["ExpiresAt"].(time.Time).Unix() < time.Now().Unix() {
		return 0, domain.TokenExpired
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid subject")
	}

	return id, nil
}

func (a *Auth) RefreshTokens(refreshToken string) (string, string, error) {
	userId, exp, err := a.usersRepo.CredentialsByRefresh(refreshToken)
	if err != nil {
		return "", "", err
	}

	if exp.Unix() < time.Now().Unix() {
		return "", "", domain.TokenExpired
	}

	return a.generateTokens(userId)
}
