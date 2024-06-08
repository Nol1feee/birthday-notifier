package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Nol1feee/birthday-notifier/internal/domain"
	"github.com/Nol1feee/birthday-notifier/internal/storage"
	"github.com/Nol1feee/birthday-notifier/pkg/logger"
)

var (
	incorrectSubInputData = "Убедись, что ты указал 3 поля (subs_id, employee_id, notify_days_before) " +
		"и notify_days_before >=0 && <=90"
	incorrectUnsubInputData = "В input доступны только int значения"
)

func (h *Handler) subscribeToEveryone(ctx *gin.Context) {
	subs := new(domain.Subscription)
	logger.Debug(fmt.Sprintf("%+v", subs))

	_ = ctx.ShouldBindJSON(&subs)

	if err := subs.Validate(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": incorrectSubInputData})
		logger.Error(err.Error())
		return
	}

	if err := h.usersService.SubscribePerPerson(subs); err != nil {
		if err == storage.DuplicateSubVal {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Подобная связь была создана ранее"})
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		logger.Error(err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "Вы успешно подписались")
}

func (h *Handler) unsubscribeFromPerson(ctx *gin.Context) {
	subs := new(domain.Subscription)
	logger.Debug(fmt.Sprintf("%+v", subs))

	_ = ctx.ShouldBindJSON(&subs)

	if err := subs.Validate(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": incorrectUnsubInputData})
		logger.Error(err.Error())
		return
	}

	if err := h.usersService.UnsubscribeFromPerson(subs); err != nil {
		if err == storage.IdSubsDoesntExists {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Подобная связь была создана ранее или айди не существуют"})
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		logger.Error(err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "Вы успешно отписались")

}
