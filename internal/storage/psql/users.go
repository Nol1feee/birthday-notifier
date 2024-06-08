package psql

import (
	"database/sql"
	"fmt"

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

func (u *Users) GetAllUsers() ([]domain.User, error) {
	//2 запроса к БД, но append не перевыделяет память
	var count int
	err := u.db.QueryRow("SELECT COUNT(*) FROM employees").Scan(&count)
	if err != nil {
		return nil, err
	}

	rows, err := u.db.Query("SELECT first_name, last_name, email, birth_date from employees")
	if err != nil {
		return nil, err
	}

	users := make([]domain.User, 0, count)

	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.FirstName, &user.LastName, &user.Email, &user.Birthdate)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *Users) GetAllBirthdayPeople(date string) ([]domain.User, error) {
	//2 запроса к БД, но append не перевыделяет память
	var count int
	err := u.db.QueryRow("SELECT COUNT(*) FROM employees WHERE TO_CHAR(birth_date, 'MM-DD') = $1", date).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("getAllBirhdayPeople | query count | error - %s", err.Error())
	}

	rows, err := u.db.Query("SELECT first_name, last_name, email FROM employees WHERE TO_CHAR(birth_date, 'MM-DD') = $1", date)
	if err != nil {
		return nil, fmt.Errorf("getAllBirhdayPeople | query birthday | error - %s", err.Error())
	}

	birthdayPeoples := make([]domain.User, 0, count)

	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.FirstName, &user.LastName, &user.Email); err != nil {
			return nil, fmt.Errorf("getAllBirhdayPeople | rows.next | error - %s", err.Error())
		}
		birthdayPeoples = append(birthdayPeoples, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getAllBirhdayPeople | rows.err | error - %s", err.Error())
	}

	return birthdayPeoples, nil
}
