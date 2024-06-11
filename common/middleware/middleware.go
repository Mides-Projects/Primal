package middleware

import (
	"net/http"
)

func HandleAuth(w http.ResponseWriter, r *http.Request) bool {
	apiKey := r.Header.Get("X-API-KEY")
	if apiKey == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	if apiKey != "supersecret" {
		http.Error(w, "Forbidden", http.StatusForbidden)

		return false
	}

	return true
}
