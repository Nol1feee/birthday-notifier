package service

import (
	"github.com/Nol1feee/birthday-notifier/internal/domain"
	"github.com/Nol1feee/birthday-notifier/internal/storage"
)

type Users struct {
	usersStorage storage.Users
}

func NewUsers(usersStorage storage.Users) *Users {
	return &Users{usersStorage: usersStorage}
}

func (u *Users) CreateUser(user *domain.User) error {
	return u.usersStorage.CreateUser(*user)
}

func (u *Users) DeleteUser(email domain.Email) error {
	return u.usersStorage.DeleteUser(email)
}

func (u *Users) GetAllUsers() ([]domain.User, error) {
	return u.usersStorage.GetAllUsers()
}

func (u *Users) SubscribePerPerson(subInfo *domain.Subscription) error {
	return u.usersStorage.SubscribePerPerson(subInfo)
}

func (u *Users) UnsubscribeFromPerson(subInfo *domain.Subscription) error {
	return u.usersStorage.UnsubscribeFromPerson(subInfo)
}
