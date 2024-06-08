package rest

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/Nol1feee/birthday-notifier/internal/service"
)

var (
	incorrectInputData = fmt.Sprintf(
		"Проверяющий из rutube, данные некорректны. Требования следующие:<br>" +
			"1. дата рождения соответствует формату 2002-10-31 (год-месяц-день) и ты старше 16 лет.<br>" +
			"2. имя и фамилия содержат от 2 до 25 символов.<br>" +
			"3. поле email валидируется, так что сервис ожидает стандартный формат email'a!<br>",
	)
	userCreated        = "Сотрудник успешно добавлен в БД!"
	userDeleted        = "Сотрудник успешно удален!"
	incorrectEmail     = "Некорректный формал email адреса"
	duplicateEmail     = errors.New("pq: duplicate key value violates unique constraint \"employees_email_key\"")
	duplicateEmailResp = "Пользователь с таким email'ом уже существует, попробуй другой адрес почты."
)

type Handler struct {
	usersService *service.Users
}

func NewHandler(usersService *service.Users) *Handler {
	return &Handler{usersService: usersService}
}

func (h *Handler) InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode) //todo fix

	r := gin.New()
	r.Use(loggerMiddleware())

	api := r.Group("/api")
	{
		api.POST("/create-employee", h.createUser)
		api.DELETE("/delete-employee/:email", h.deleteUser)
		api.GET("/all-employees", h.getAllUsers)
		//внутренняя логика - поздравление с ДР именинника

		//api.POST("/subscription-all")        //с настройкой кол-ва дней до др
		//api.POST("/subscription-per-person") //с настройкой кол-ва дней до др | 1 запрос = 1 подписка
		//
		//api.POST("/unsubscribe") //в зависимости от параметра, т.е. либо 1 отдельный человек, либо все
	}
	return r
}

//EMPLOYEES -> id, firstName, lastName, email, date of birthd,created_at
//
