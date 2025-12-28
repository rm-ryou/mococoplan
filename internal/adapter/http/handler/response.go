package handler

import (
	"encoding/json"
	"net/http"
)

func writeJson(w http.ResponseWriter, statusCode int, item any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if item == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(item); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func writeErr(w http.ResponseWriter, statusCode int, err error) {
	type errResponse struct {
		error   string
		message string
	}

	res := errResponse{
		error:   http.StatusText(statusCode),
		message: err.Error(),
	}

	writeJson(w, statusCode, res)
}
