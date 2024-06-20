package storage

import (
	"errors"

	"github.com/Nol1feee/birthday-notifier/internal/domain"
)

const (
	dbName = "em"
)

var (
	EmailDoesntExits   = errors.New("запрашиваемый email для удаления не существует")
	DuplicateSubVal    = errors.New("input = дубль")
	IdSubsDoesntExists = errors.New("такие айди не сущесвуют")
)

type Users interface {
	CreateUser(user domain.User) error
	DeleteUser(email domain.Email) error
	GetAllUsers() ([]domain.User, error)
	GetAllBirthdayPeople(date string) ([]domain.User, error)
	SubscribePerPerson(subsInfo *domain.Subscription) error
	UnsubscribeFromPerson(subInfo *domain.Subscription) error
	GetBirthdayNotifications() (map[string][]*domain.User, error)
}
