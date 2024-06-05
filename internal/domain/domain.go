package domain

type (
	User struct {
		FirstName string `json:"first_name" validate:"required,lte=25,gte=2"`
		LastName  string `json:"last_name" validate:"required,lte=25,gte=2"`
		Email     string `json:"email" validate:"required,email"`
		Birthdate string `json:"birthdate" validate:"required,birthdate"`
	}

	Email struct {
		Email string `json:"email" validate:"required,email"`
	}
)
