package domain

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("birthdate", birthdateValidation)
}

func (s User) Validate() error {
	return validate.Struct(s)
}

func (s Email) Validate() error {
	return validate.Struct(s)
}

// сотрудник компании должен быть старше 16 лет, -> валидируем
func birthdateValidation(fl validator.FieldLevel) bool {
	birthdateStr := fl.Field().String()
	layout := "2006-01-02"
	birthdate, err := time.Parse(layout, birthdateStr)
	if err != nil {
		return false
	}

	today := time.Now() //
	xYearsAgo := today.AddDate(-16, 0, 0)
	return birthdate.Before(xYearsAgo)
}
