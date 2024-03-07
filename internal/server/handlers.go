package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Richtermnd/TgLogin/internal/storage"
	"github.com/Richtermnd/TgLogin/pkg/tglogin"
)

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	userData := tglogin.FromJSON(r.Body)
	tglogin.SetCookie(w, userData)

	fmt.Printf("userData: %v\n", userData)
	// If auth data is invalid, send user to login page
	if !s.service.IsAuthentificated(r.Context(), userData) {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	// Check user exist
	_, err := s.service.User(r.Context(), userData)
	// If user exist fix login and send OK
	if err == nil {
		s.service.Login(r.Context(), userData)
		w.WriteHeader(http.StatusOK)
		return
	}

	// If error is not "not found" send InternalServerError
	if !errors.Is(err, storage.ErrNotFound) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// If user not found register him
	err = s.service.RegisterUser(r.Context(), userData)
	// If errors send InternalServerError else OK
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) user(w http.ResponseWriter, r *http.Request) {
	userData, err := tglogin.FromCookie(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	user, err := s.service.User(r.Context(), userData)
	if err != nil {
		var code int
		if errors.Is(err, storage.ErrNotFound) {
			code = http.StatusNotFound
		} else {
			code = http.StatusInternalServerError
		}
		http.Redirect(w, r, "/login", code)
	}
	if err := encodeUser(w, user); err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	tglogin.DeleteCookie(w)
	http.Redirect(w, r, "/login", http.StatusOK)
}

// I know about http.FileServer
func (s *Server) loginPage(w http.ResponseWriter, r *http.Request) {
	w.Write(loginPage)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) userPage(w http.ResponseWriter, r *http.Request) {
	w.Write(userPage)
	w.WriteHeader(http.StatusOK)
}
