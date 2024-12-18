package utils

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ExtractEventID(res http.ResponseWriter, req *http.Request) (int, bool) {
	idStr := chi.URLParam(req, "id")
	id, err := strconv.Atoi(idStr)

	ok := true

	if err != nil {
		http.Error(res, "Invalid event ID", http.StatusBadRequest)
		ok = false
	}

	return id, ok
}
