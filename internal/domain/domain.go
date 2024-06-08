package domain

type (
	User struct {
		Id        int    `json:"id"`
		FirstName string `json:"first_name" validate:"required,lte=25,gte=2"`
		LastName  string `json:"last_name" validate:"required,lte=25,gte=2"`
		Email     string `json:"email" validate:"required,email"`
		Birthdate string `json:"birthdate" validate:"required,birthdate"`
	}

	Email struct {
		Email string `json:"email" validate:"required,email"`
	}

	Subscription struct {
		SubsId           int `json:"subs_id" validate:"required"`
		EmployeeId       int `json:"employee_id" validate:"required"`
		NotifyDaysBefore int `json:"notify_days_before" validate:"min=0,max=90"`
	}
)
