package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/anatolyi0311/go-musthave-diploma-tpl/pkg/helpers"
	"github.com/anatolyi0311/go-musthave-diploma-tpl/pkg/user"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	var req user.UserRequest

	isHashed := r.Header.Get("X-Password-Format") == "sha256"

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "Login and password are required", http.StatusBadRequest)
		return
	}

	userID, err := h.Storage.DBStorage.GetUserID(req, isHashed)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid login/password", http.StatusUnauthorized)
		return
	}

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	cookie, err := helpers.GetAuthCookie(userID, req.Login, h.JwtKey)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}