package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

//
//  Middleware
//
func Gzip(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		c.Response.Header().Set("Vary", "Accept-Encoding")

		// do nothing if the browser does not accept gzip
		if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
			return h(c)
		}

		// headers
		c.Response.Header().Set("Content-Encoding", "gzip")

		// gzip writer
		gz := gzip.NewWriter(c.Response.Writer)
		defer gz.Close()

		// gzip response writer
		gzw := gzipResponseWriter{Writer: gz, ResponseWriter: c.Response.Writer}

		// change context response writer
		c.Response.Writer = gzw

		// return
		return h(c)
	}
}

//
//  ResponseWriter Write function
//
func (grw gzipResponseWriter) Write(p []byte) (int, error) {
	if len(grw.Header().Get("Content-Type")) == 0 {
		grw.Header().Set("Content-Type", http.DetectContentType(p))
	}

	return grw.Writer.Write(p)
}
