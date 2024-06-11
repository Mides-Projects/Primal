package middleware

import (
	"github.com/holypvp/primal/common"
	"net/http"
)

func HandleAuth(w http.ResponseWriter, r *http.Request) bool {
	apiKey := r.Header.Get("API-Key")
	if apiKey == common.APIKey {
		return true
	}

	if apiKey == "" {
		http.Error(w, "Forbidden", http.StatusForbidden)

		return false
	}

	http.Error(w, "Unauthorized", http.StatusUnauthorized)

	return false
}
