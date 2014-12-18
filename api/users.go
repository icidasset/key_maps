package api

import (
  "github.com/dgrijalva/jwt-go"
  _ "github.com/go-martini/martini"
  "github.com/icidasset/key-maps/db"
  _ "github.com/lib/pq"
  "github.com/martini-contrib/render"
  "golang.org/x/crypto/bcrypt"
  "time"
)


type User struct {
  Id int
  Email string
  EncryptedPassword string    `db:"encrypted_password"`
  CreatedAt time.Time         `db:"created_at"`
  UpdatedAt time.Time         `db:"updated_at"`
}


type UserAuth struct {
  Email string                `form:"email" binding:"required"`
  Password string             `form:"password" binding:"required"`
}


type UserAuthFormData struct {
  User UserAuth               `form:"user" binding:"required"`
}


type UserPublic struct {
  Token string                `json:"token"`
}



//
//  Routes
//
func Users__Create(ufd UserAuthFormData, r render.Render) {
  query := "INSERT INTO users (email, encrypted_password, created_at, updated_at) VALUES (:email, :encrypted_password, :created_at, :updated_at)"

  // make new user
  encryped_password, _ := bcrypt.GenerateFromPassword(
    []byte(ufd.User.Password),
    bcrypt.DefaultCost,
  )

  now := time.Now()

  new_user := User{
    Email: ufd.User.Email,
    EncryptedPassword: string(encryped_password),
    CreatedAt: now,
    UpdatedAt: now,
  }

  // execute query
  _, err := db.Inst().NamedExec(query, new_user)

  // if error
  if err != nil {
    r.JSON(500, err.Error())
    return
  }

  // render user as json
  user := User{}

  db.Inst().Get(
    &user,
    "SELECT * FROM users WHERE email = $1 LIMIT 1",
    new_user.Email,
  )

  token := generate_new_token(&user)
  user_public := UserPublic{ Token: token }

  r.JSON(200, map[string]UserPublic{ "user": user_public })
}


func Users__Authenticate(ufd UserAuthFormData, r render.Render) {
  user := User{}

  db.Inst().Get(
    &user,
    "SELECT * FROM users WHERE email = $1 LIMIT 1",
    ufd.User.Email,
  )

  if user.Email == "" {
    // {err} user doesn't exist
    r.JSON(500, map[string]string{ "error" : "User not found." })
  }

  bcrypt_check_err := bcrypt.CompareHashAndPassword(
    []byte(user.EncryptedPassword),
    []byte(ufd.User.Password),
  )

  if bcrypt_check_err != nil {
    // {err} invalid password
    r.JSON(500, map[string]string{ "error" : "Invalid password." })
  }

  token := generate_new_token(&user)
  user_public := UserPublic{ Token: token }

  r.JSON(200, map[string]UserPublic{ "user": user_public })
}



//
//  Helpers
//
func generate_new_token(user *User) string {
  token := jwt.New(jwt.GetSigningMethod("HS256"))
  token.Claims["user_id"] = user.Id
  token.Claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
  token_string, _ := token.SignedString([]byte("TODO - SECRET KEY"))

  return token_string
}
