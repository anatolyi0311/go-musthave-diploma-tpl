package handler

import (
	"encoding/json"
	"net/http"

	"github.com/anatolyi0311/go-musthave-diploma-tpl/pkg/helpers"
	"github.com/anatolyi0311/go-musthave-diploma-tpl/pkg/user"
)

func (h *Handler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {

	isHashed := r.Header.Get("X-Password-Format") == "sha256"

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	var req user.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "Login and password are required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	userIsRegistred, err := h.Storage.DBStorage.UserIsRegistred(req.Login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if userIsRegistred {
		w.WriteHeader(http.StatusConflict)
		return
	}

	userID, err := h.Storage.DBStorage.RegisterUserWithBalance(req, isHashed)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie, err := helpers.GetAuthCookie(userID, req.Login, h.JwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})

}