package storage

import (
	"errors"

	"github.com/Nol1feee/birthday-notifier/internal/domain"
)

const (
	dbName = "em"
)

var (
	EmailDoesntExits = errors.New("запрашиваемый email для удаления не существует")
)

type Users interface {
	CreateUser(user domain.User) error
	DeleteUser(email domain.Email) error
	GetAllUsers() ([]domain.User, error)
	GetAllBirthdayPeople(date string) ([]domain.User, error)
}
