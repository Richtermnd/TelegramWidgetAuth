package server

import (
	"fmt"
	"net/http"

	"github.com/Richtermnd/TgLogin/pkg/tglogin"
)

func (s *Server) loginRequiredMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userData, err := tglogin.FromCookie(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}
		if !s.service.IsAuthentificated(r.Context(), userData) {
			fmt.Println("Middleware block")
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
		}
		next.ServeHTTP(w, r)
	}
}
