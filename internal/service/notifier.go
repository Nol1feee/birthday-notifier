package service

import (
	"context"
	"fmt"
	"time"

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
	upcomingBirthdays    = "У вас есть предстоящие дни рождения:\n<ul>"
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

func (n *Notifier) sendEmail(email string, subject string, body string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", n.cfg.From)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)

	m.SetBody("text/html", body)

	if err := n.dialer.DialAndSend(m); err != nil {
		logger.Error("email errors", zap.Error(err))
		return err
	}

	return nil
}

func (n *Notifier) CongratulateAll(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute) //условно
	defer cancel()

	todayDate, err := today()
	if err != nil {
		logger.Error("get today")
		return err
	}

	users, err := n.usersStorage.GetAllBirthdayPeople(ctx, todayDate)
	logger.Debug("get all today's birthday",
		zap.String("peoples - ", fmt.Sprintf("%v", users)))
	if err != nil {
		logger.Error("get all today birthday guys")
		return err
	}

	for _, v := range users {
		body := fmt.Sprintf(bodyHappyBirthday, v.FirstName, v.LastName)
		err := n.sendEmail(v.Email, subjectHappyBirthday, body)
		if err != nil {
			logger.Error("Ошибка при отправке email", zap.String(v.FirstName+v.LastName, "не получила письмо"))
		}
	}

	return nil
}

func (n *Notifier) NotifyingUpcomingBirthdays(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute) //условно
	defer cancel()

	notifications, err := n.usersStorage.GetBirthdayNotifications(ctx)
	if err != nil {
		logger.Error("Error fetching birthday notifications", zap.Error(err))
		return err
	}

	for subscriberEmail, birthdayPeople := range notifications {
		body := generateNotificationBody(birthdayPeople)
		subject := "Upcoming Birthdays Notification"

		if err := n.sendEmail(subscriberEmail, subject, body); err != nil {
			logger.Error("Failed to send notification email",
				zap.String("subscriber_email", subscriberEmail),
				zap.Error(err))
			return err
		}
		logger.Info("Notification email sent successfully",
			zap.String("subscriber_email", subscriberEmail))
	}
	return err
}

func generateNotificationBody(birthdayPeople []*domain.User) string {
	body := upcomingBirthdays
	for _, person := range birthdayPeople {
		body += fmt.Sprintf("<li>%s %s - %s</li>", person.FirstName, person.LastName, person.Birthdate)
	}
	body += "</ul>"
	return body
}
