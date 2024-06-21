package psql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/Nol1feee/birthday-notifier/internal/domain"
	"github.com/Nol1feee/birthday-notifier/internal/storage"
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{db: db}
}

func (u *Users) CreateUser(ctx context.Context, user domain.User) error {
	_, err := u.db.Exec("INSERT INTO employees (first_name, last_name, email, birth_date) values ($1, $2, $3, $4)",
		user.FirstName, user.LastName, user.Email, user.Birthdate)

	return err
}

func (u *Users) DeleteUser(ctx context.Context, email domain.Email) error {
	err := u.db.QueryRow("SELECT email from employees WHERE email=$1)", email.Email).Err()
	if err == nil {
		_, err = u.db.Exec("DELETE FROM employees WHERE email=$1", email.Email)
		return err
	}

	return storage.EmailDoesntExits
}

func (u *Users) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	//2 запроса к БД, но append не перевыделяет память
	var count int
	err := u.db.QueryRow("SELECT COUNT(*) FROM employees").Scan(&count)
	if err != nil {
		return nil, err
	}

	rows, err := u.db.Query("SELECT id, first_name, last_name, email, birth_date from employees")
	if err != nil {
		return nil, err
	}

	users := make([]domain.User, 0, count)

	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Birthdate)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *Users) GetAllBirthdayPeople(ctx context.Context, date string) ([]domain.User, error) {
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

func (u *Users) SubscribePerPerson(ctx context.Context, info *domain.Subscription) error {
	_, err := u.db.Exec("INSERT INTO subscriptions(subscriber_id, employee_id, notify_days_before) values ($1, $2, $3)",
		info.SubsId, info.EmployeeId, info.NotifyDaysBefore)

	//проверка, чо ошибка = дубль
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return storage.DuplicateSubVal
			}
		}
	}

	return err
}

func (u *Users) UnsubscribeFromPerson(ctx context.Context, subInfo *domain.Subscription) error {
	var id int
	err := u.db.QueryRow("SELECT id from subscriptions WHERE subscriber_id=$1 AND employee_id=$2", subInfo.SubsId, subInfo.EmployeeId).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return storage.IdSubsDoesntExists
		}
		return err
	}

	_, err = u.db.Exec("DELETE FROM subscriptions WHERE subscriber_id=$1 AND employee_id=$2", subInfo.SubsId, subInfo.EmployeeId)
	return err
}

func (u *Users) GetBirthdayNotifications(ctx context.Context) (map[string][]*domain.User, error) {
	query := `
        SELECT 
            e1.email AS subscriber_email,
            e2.id AS birthday_person_id,
            e2.first_name AS birthday_person_first_name,
            e2.last_name AS birthday_person_last_name,
            e2.email AS birthday_person_email,
            e2.birth_date AS birthday_person_birth_date
        FROM 
            subscriptions s
        JOIN 
            employees e1 ON s.subscriber_id = e1.id
        JOIN 
            employees e2 ON s.employee_id = e2.id
        WHERE 
            TO_CHAR(e2.birth_date - INTERVAL '1 day' * s.notify_days_before, 'MM-DD') = TO_CHAR(CURRENT_DATE, 'MM-DD')
        ORDER BY 
            e1.email;
    `

	rows, err := u.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	notifications := make(map[string][]*domain.User)

	for rows.Next() {
		var subscriberEmail string
		birthdayPerson := &domain.User{}

		if err := rows.Scan(&subscriberEmail, &birthdayPerson.Id, &birthdayPerson.FirstName, &birthdayPerson.LastName, &birthdayPerson.Email, &birthdayPerson.Birthdate); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		notifications[subscriberEmail] = append(notifications[subscriberEmail], birthdayPerson)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	return notifications, nil
}
