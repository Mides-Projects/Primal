package middleware

import (
	"net/http"
)

func HandleAuth(w http.ResponseWriter, r *http.Request) bool {
	apiKey := r.Header.Get("X-API-KEY")
	if apiKey == "superuser" {
		return true
	}

	if apiKey == "" {
		http.Error(w, "Forbidden", http.StatusForbidden)

		return false
	}

	http.Error(w, "Unauthorized", http.StatusUnauthorized)

	return false
}
