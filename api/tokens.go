package api

import (
  "github.com/dgrijalva/jwt-go"
  "os"
  "time"
)


var SECRET_KEY []byte = []byte(
  os.Getenv("SECRET_KEY"),
)


func GenerateToken(user *User) string {
  token := jwt.New(jwt.GetSigningMethod("HS256"))
  token.Claims["user_id"] = int(user.Id)
  token.Claims["exp"] = time.Now().Add(time.Hour * 24 * 365).Unix()
  token_string, _ := token.SignedString(SECRET_KEY)

  return token_string
}


func ParseToken(t string) *jwt.Token {
  token, _ := jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
    return SECRET_KEY, nil
  })

  return token
}


func VerifyToken(t string) bool {
  token := ParseToken(t)
  return token.Valid;
}
