package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Richtermnd/TgLogin/internal/config"
	"github.com/Richtermnd/TgLogin/internal/domain"
	"github.com/Richtermnd/TgLogin/pkg/tglogin"
)

type Service interface {
	User(ctx context.Context, userData tglogin.TelegramUserData) (domain.User, error)
	RegisterUser(ctx context.Context, userData tglogin.TelegramUserData) error
	Login(ctx context.Context, userData tglogin.TelegramUserData) error
	IsAuthentificated(ctx context.Context, userData tglogin.TelegramUserData) bool
}

type Server struct {
	service Service
	server  *http.Server
}

func New(service Service) *Server {
	addr := fmt.Sprintf("localhost:%d", config.Config().Port)
	mux := http.NewServeMux()

	httpServer := http.Server{Addr: addr, Handler: mux}
	s := &Server{
		service: service,
		server:  &httpServer,
	}
	registerHandlers(mux, s)
	return s
}

func (s *Server) Run() {
	log.Printf("Server started on %s", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil {
		log.Println("Server stopped")
	}
}

func (s *Server) Shutdown() {
	log.Println("Shutdown server")
	s.server.Shutdown(context.Background())
}

func registerHandlers(mux *http.ServeMux, s *Server) {
	uploadPages()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusFound)
	})
	mux.HandleFunc("GET /login", s.loginPage)
	mux.HandleFunc("POST /api/login", s.login)
	mux.HandleFunc("GET /logout", s.logout)
	mux.HandleFunc("GET /user", s.loginRequiredMiddleware(s.userPage))
	mux.HandleFunc("GET /api/user", s.loginRequiredMiddleware(s.user))
}
