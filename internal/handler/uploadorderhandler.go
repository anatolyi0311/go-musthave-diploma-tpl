package handler

import (
	"database/sql"
	"io"
	"net/http"
	"strconv"

	"github.com/anatolyi0311/go-musthave-diploma-tpl/pkg/user"
	"github.com/gymgle/exercism/go/luhn"
)

func (h *Handler) UploadOrderHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(user.UserIDKey).(int)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	orderNumber := string(body)
	if orderNumber == "" {
		http.Error(w, "Empty order number", http.StatusBadRequest)
	}
	_, err = strconv.ParseInt(orderNumber, 10, 64)
	if err != nil || !luhn.Valid(orderNumber) {
		http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
		return
	}

	// Проверка существования заказа
	existingUserID, err := h.Storage.DBStorage.GetUserIDWithOrder(orderNumber)
	if err == nil {
		if existingUserID == userID {
			w.WriteHeader(http.StatusOK)
			return
		} else {
			http.Error(w, "Order already uploaded by another user", http.StatusConflict)
			return
		}
	}

	if err != sql.ErrNoRows {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Создание нового заказа
	if err = h.Storage.DBStorage.AddNewOrder(orderNumber, userID, "NEW"); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}