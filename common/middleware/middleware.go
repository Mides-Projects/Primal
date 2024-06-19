package middleware

import (
	"github.com/holypvp/primal/common"
	"github.com/labstack/echo/v4"
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

func HandleBasicAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		apiKey := c.Request().Header.Get("X-API-Key")
		if apiKey == common.APIKey {
			return next(c)
		}

		if apiKey == "" {
			return c.String(http.StatusForbidden, "Forbidden")
		}

		return c.String(http.StatusUnauthorized, "Unauthorized")
	}
}
