package respond

import (
	"WalletX/pkg/errs"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, status int, message string, err error) {
	JSON(w, status, map[string]interface{}{
		"error":   message,
		"details": err.Error(),
	})
}

func HandleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, errs.ErrInvalidPhone),
		errors.Is(err, errs.ErrUserExists),
		errors.Is(err, errs.ErrWeakPassword),
		errors.Is(err, errs.ErrValidationFailed),
		errors.Is(err, errs.ErrRequiredFields):
		JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})

	case errors.Is(err, errs.ErrUserNotFound):
		JSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})

	case errors.Is(err, errs.ErrUnauthorized):
		JSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})

	default:
		JSON(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("something went wrong: %s", err.Error()),
		})
	}
}
