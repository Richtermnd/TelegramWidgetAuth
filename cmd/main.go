package main

import (
	"log/slog"
	"os"
	"os/signal"

	"github.com/Richtermnd/TgLogin/internal/app"
	"github.com/Richtermnd/TgLogin/internal/config"
)

func main() {
	config.LoadConfig()
	log := slog.Default()

	app := app.New(log)
	app.Run()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	app.Shutdown()
}
