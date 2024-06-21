package storage

import (
	"context"
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
	CreateUser(ctx context.Context, user domain.User) error
	DeleteUser(ctx context.Context, email domain.Email) error
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	GetAllBirthdayPeople(ctx context.Context, date string) ([]domain.User, error)
	SubscribePerPerson(ctx context.Context, subsInfo *domain.Subscription) error
	UnsubscribeFromPerson(ctx context.Context, subInfo *domain.Subscription) error
	GetBirthdayNotifications(ctx context.Context) (map[string][]*domain.User, error)
}
