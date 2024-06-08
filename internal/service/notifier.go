package service

import (
	"fmt"

	"go.uber.org/zap"
	"gopkg.in/gomail.v2"

	"github.com/Nol1feee/birthday-notifier/config"
	"github.com/Nol1feee/birthday-notifier/internal/domain"
	"github.com/Nol1feee/birthday-notifier/internal/storage"
	"github.com/Nol1feee/birthday-notifier/pkg/logger"
)

var (
	subjectHappyBirthday = "happy birthday"
	bodyHappyBirthday    = "Дорогой %s %s, мы поздравляем тебя с днем рождения!"
)

type Notifier struct {
	cfg          config.Email
	usersStorage storage.Users
	dialer       *gomail.Dialer
}

func NewNotifier(cfg config.Email, usersStorage storage.Users) *Notifier {
	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.From, cfg.Pass)
	return &Notifier{cfg: cfg, usersStorage: usersStorage, dialer: dialer}
}

func (n *Notifier) sendEmail(user domain.User, subject string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", n.cfg.From)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", subject)

	m.SetBody("text/html", fmt.Sprintf(bodyHappyBirthday, user.FirstName, user.LastName))

	if err := n.dialer.DialAndSend(m); err != nil {
		logger.Error("email errors", zap.Error(err))
		return err
	}

	return nil
}

func (n *Notifier) CongratulateAll() error {
	todayDate, err := today()
	if err != nil {
		logger.Error("get today")
		return err
	}

	users, err := n.usersStorage.GetAllBirthdayPeople(todayDate)
	logger.Debug("get all today's birthday",
		zap.String("peoples - ", fmt.Sprintf("%v", users)))
	if err != nil {
		logger.Error("get all today birthday guys")
		return err
	}

	for _, v := range users {
		err := n.sendEmail(v, subjectHappyBirthday)
		if err != nil {
			logger.Error("Ошибка при отправке email", zap.String(v.FirstName+v.LastName, "не получила письмо"))
		}
	}

	return nil
}
