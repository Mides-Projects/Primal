package middleware

import (
	"github.com/holypvp/primal/common"
	"github.com/labstack/echo/v4"
	"net/http"
)

func HandleBasicAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		apiKey := c.Request().Header.Get("X-API-Key")
		if apiKey == common.APIKey {
			return next(c)
		}

		if apiKey == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "API key is required")
		}

		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}
}
