package auth

import (
	"errors"
	"fmt"
	"os"

	"go.uber.org/zap"
	"gopkg.in/gomail.v2"

	"github.com/FreshOfficeFriends/SSO/pkg/logger"
)

func SendEmail(email string, hashEmail string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", "vladrazd00@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "confirmation code")

	m.SetBody("text/html", fmt.Sprintf("confirm your email - <b>http://localhost:8080/sso/auth/confirm-email/%s</b>", hashEmail))
	msg := fmt.Sprintf("confirm your email - <b>http://localhost:9000/sso/auth/confirm-email/%s</b>", hashEmail)
	logger.Debug(msg)

	d := gomail.NewDialer("smtp.gmail.com", 587, "vladrazd00@gmail.com", os.Getenv("GMAIL_PASS"))

	if err := d.DialAndSend(m); err != nil {
		logger.Error("email errors", zap.Error(err))
		return errors.New("не удалось отправить email")
	}
	return nil
}
