package rest

import (
	"github.com/gin-gonic/gin"

	"github.com/Nol1feee/birthday-notifier/internal/service"
)

type Handler struct {
	service.Users
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(loggerMiddleware())

	api := r.Group("/api")
	{
		api.POST("/create-employee", h.createUser)
		api.DELETE("/delete-employee/:id")
		api.GET("/all-employees")
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
