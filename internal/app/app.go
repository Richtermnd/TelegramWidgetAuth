package app

import (
	"log/slog"

	"github.com/Richtermnd/TgLogin/internal/server"
	"github.com/Richtermnd/TgLogin/internal/service"
	"github.com/Richtermnd/TgLogin/internal/storage/sqlite"
)

type App struct {
	log    *slog.Logger
	server *server.Server
}

func New(log *slog.Logger) *App {
	storage := sqlite.New()
	service := service.New(log, storage)
	server := server.New(service)
	return &App{log: log, server: server}
}

func (a *App) Run() {
	go a.server.Run()
}

func (a *App) Shutdown() {
	a.server.Shutdown()
}
