package service

import (
	"context"

	"github.com/Nol1feee/birthday-notifier/internal/domain"
	"github.com/Nol1feee/birthday-notifier/internal/storage"
)

type Users struct {
	usersStorage storage.Users
}

func NewUsers(usersStorage storage.Users) *Users {
	return &Users{usersStorage: usersStorage}
}

func (u *Users) CreateUser(ctx context.Context, user *domain.User) error {
	return u.usersStorage.CreateUser(ctx, *user)
}

func (u *Users) DeleteUser(ctx context.Context, email domain.Email) error {
	return u.usersStorage.DeleteUser(ctx, email)
}

func (u *Users) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	return u.usersStorage.GetAllUsers(ctx)
}

func (u *Users) SubscribePerPerson(ctx context.Context, subInfo *domain.Subscription) error {
	return u.usersStorage.SubscribePerPerson(ctx, subInfo)
}

func (u *Users) UnsubscribeFromPerson(ctx context.Context, subInfo *domain.Subscription) error {
	return u.usersStorage.UnsubscribeFromPerson(ctx, subInfo)
}
