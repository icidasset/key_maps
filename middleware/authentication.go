package middleware

import (
	"net/http"
	"strings"

	"github.com/icidasset/key-maps-api/api"
	"github.com/labstack/echo"
)

func MustBeAuthenticated(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {

		if c.Request().Method == "OPTIONS" {
			return h(c)
		}

		auth_header := c.Request().Header.Get("Authorization")

		if strings.Contains(auth_header, "Bearer") {
			t := strings.Split(auth_header, "Bearer ")[1]
			token, err := api.ParseToken(t)
			is_valid_token := false

			if err == nil && token.Valid {
				is_valid_token = true
			}

			if !is_valid_token {
				return c.HTML(http.StatusUnauthorized, "Forbidden")

			} else {
				id := int(token.Claims["user_id"].(float64))
				c.Set("user", api.User{Id: id})
				return h(c)
			}

		} else {
			return c.HTML(http.StatusUnauthorized, "Forbidden")

		}

	}
}
