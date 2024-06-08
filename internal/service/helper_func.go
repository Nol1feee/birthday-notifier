package service

import (
	"fmt"
	"time"
)

func today() (string, error) {
	location, err := time.LoadLocation("Europe/Moscow")

	if err != nil {
		return "", fmt.Errorf("ошибка загрузки временной зоны Москвы: %v", err)
	}

	now := time.Now().In(location)

	return now.Format("01-02"), nil
}
