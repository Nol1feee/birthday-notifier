package psql

import (
	"database/sql"

	"github.com/Nol1feee/birthday-notifier/internal/domain"
	"github.com/Nol1feee/birthday-notifier/internal/storage"
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{db: db}
}

func (u *Users) CreateUser(user domain.User) error {
	_, err := u.db.Exec("INSERT INTO employees (first_name, last_name, email, birth_date) values ($1, $2, $3, $4)",
		user.FirstName, user.LastName, user.Email, user.Birthdate)

	return err
}

func (u *Users) DeleteUser(email domain.Email) error {
	err := u.db.QueryRow("SELECT email from employees WHERE email=$1)", email.Email).Err()
	if err == nil {
		_, err = u.db.Exec("DELETE FROM employees WHERE email=$1", email.Email)
		return err
	}

	return storage.EmailDoesntExits
}
