package auth

import (
	"errors"
	"fmt"
	"os"

	"go.uber.org/zap"
	"gopkg.in/gomail.v2"

	"github.com/FreshOfficeFriends/SSO/pkg/logger"
)

//var confirmEmail = fmt.Sprintf("confirm your email - <b>http://localhost:8080/sso/auth/confirm-email/%s</b>")
//
//type EmailInfo struct {
//	To      string
//	Subject string
//	Body    string
//}
//
//func (e *EmailInfo) Send(info EmailInfo) error {
//	m := gomail.NewMessage()
//
//	m.SetHeader("From", e.email.From)
//	m.SetHeader("To", info.To)
//	m.SetHeader("Subject", info.Subject)
//
//	m.SetBody("text/html", info.Body)
//	logger.Debug(info.Body)
//
//	d := gomail.NewDialer(e.email.Host, e.email.Port, e.email.From, e.email.Pass)
//
//	if err := d.DialAndSend(m); err != nil {
//		logger.Error("email errors", zap.Error(err))
//		return errors.New("не удалось отправить email")
//	}
//	return nil
//}

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

func SendEmailRecoverPass(email string, uuid string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", "vladrazd00@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "recover password")

	//uuid
	m.SetBody("text/html", fmt.Sprintf("recover your password - <b>http://localhost:8080/sso/auth/reset-password/%s</b>", uuid))
	msg := fmt.Sprintf("recover your password - <b>http://localhost:8080/sso/auth/reset-password/%s</b>", uuid)
	logger.Debug(msg)

	d := gomail.NewDialer("smtp.gmail.com", 587, "vladrazd00@gmail.com", os.Getenv("GMAIL_PASS"))

	if err := d.DialAndSend(m); err != nil {
		logger.Error("email errors", zap.Error(err))
		return errors.New("не удалось отправить email")
	}
	return nil
}
