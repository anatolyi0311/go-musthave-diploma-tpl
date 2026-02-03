package handler

import (
	"encoding/json"
	"net/http"

	"github.com/anatolyi0311/go-musthave-diploma-tpl/pkg/helpers"
	"github.com/anatolyi0311/go-musthave-diploma-tpl/pkg/user"
)

func (h *Handler) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(user.UserIDKey).(int)

	w.Header().Set("Content-Type", "application/json")

	var userBalance user.Balance
	row := h.Storage.DBStorage.GetUserBalance(userID)

	userBalance, err := helpers.GetUserBalance(row)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userBalance)
}
