package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Nol1feee/birthday-notifier/internal/domain"
	"github.com/Nol1feee/birthday-notifier/internal/storage"
	"github.com/Nol1feee/birthday-notifier/pkg/logger"
)

func (h *Handler) createUser(ctx *gin.Context) {
	user := new(domain.User)

	_ = ctx.ShouldBindJSON(&user)

	if err := user.Validate(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": incorrectInputData})
		logger.Error(err.Error())
		return
	}

	//доп. обработка ошибки, т.к. быть может такой email существует
	if err := h.usersService.CreateUser(user); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		logger.Error(err.Error())
		return
	}

	ctx.JSON(http.StatusOK, userCreated)
}

func (h *Handler) deleteUser(ctx *gin.Context) {
	email := domain.Email{Email: ctx.Param("email")}

	if err := email.Validate(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": incorrectEmail})
		logger.Error(err.Error())
		return
	}

	if err := h.usersService.DeleteUser(email); err != nil {
		fmt.Println(err, storage.EmailDoesntExits)
		if errors.Is(err, storage.EmailDoesntExits) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, storage.EmailDoesntExits.Error())
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		logger.Error(err.Error())
		return
	}

	ctx.JSON(http.StatusOK, userDeleted)

}

// 2 min, little chill
func (h *Handler) getAllUsers(ctx *gin.Context) {

}
