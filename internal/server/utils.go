package server

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/Richtermnd/TgLogin/internal/domain"
)

var (
	loginPage []byte
	userPage  []byte
)

func encodeUser(w http.ResponseWriter, user domain.User) error {
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return err
}

func uploadPages() {
	f, _ := os.Open("./frontend/login.html")
	defer f.Close()
	loginPage, _ = io.ReadAll(f)

	f, _ = os.Open("./frontend/user.html")
	defer f.Close()
	userPage, _ = io.ReadAll(f)
}
