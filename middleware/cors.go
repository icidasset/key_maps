package middleware

import "github.com/labstack/echo"

func Cors(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		c.Response.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, PATCH, DELETE")
		c.Response.Header().Set("Access-Control-Allow-Origin", "*")

		return h(c)
	}
}
