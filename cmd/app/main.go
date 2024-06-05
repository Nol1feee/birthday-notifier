package main

import (
	"github.com/Nol1feee/birthday-notifier/config"
	"github.com/Nol1feee/birthday-notifier/internal/app"
)

func main() {
	app.Run(config.New())
}
