package api

import (
  "github.com/gocraft/web"
  "net/http"
  "strings"
)


func (c *Context) MustBeAuthenticated(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
  auth_header := req.Header.Get("Authorization")

  if strings.Contains(auth_header, "Bearer") {
    t := strings.Split(auth_header, "Bearer ")[1]
    token := ParseToken(t)

    if !token.Valid {
      http.Error(rw, "Forbidden", http.StatusUnauthorized)
    } else {
      id := int(token.Claims["user_id"].(float64))
      c.User = User{ Id: id }
      next(rw, req)
    }

  } else {
    http.Error(rw, "Forbidden", http.StatusUnauthorized)

  }
}
