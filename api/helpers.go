package api

import (
  "github.com/dgrijalva/jwt-go"
  "time"
)


func GenerateToken(user *User) string {
  token := jwt.New(jwt.GetSigningMethod("HS256"))
  token.Claims["user_id"] = user.Id
  token.Claims["exp"] = time.Now().Add(time.Minute * 60).Unix()
  token_string, _ := token.SignedString([]byte("TODO - SECRET KEY"))

  return token_string
}


func VerifyToken(t string) bool {
  token, _ := jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
    return []byte("TODO - SECRET KEY"), nil
  })

  return token.Valid;
}
