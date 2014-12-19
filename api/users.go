package api

import (
  "github.com/icidasset/key-maps/db"
  "github.com/martini-contrib/render"
  "golang.org/x/crypto/bcrypt"
  "net/http"
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
  Email string                `form:"email"`
  Password string             `form:"password"`
}


type UserAuthFormData struct {
  User UserAuth               `form:"user"`
}


type UserPublic struct {
  Token string                `json:"token"`
}



//
//  -> CREATE
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
    r.JSON(500, map[string]string{ "error": err.Error() })
    return
  }

  // render user as json
  user := User{}

  db.Inst().Get(
    &user,
    "SELECT * FROM users WHERE email = $1 LIMIT 1",
    new_user.Email,
  )

  token := GenerateToken(&user)
  user_public := UserPublic{ Token: token }

  r.JSON(200, map[string]UserPublic{ "user": user_public })
}



//
//  -> AUTHENTICATE
//
func Users__Authenticate(ufd UserAuthFormData, r render.Render) {
  user := User{}

  db.Inst().Get(
    &user,
    "SELECT * FROM users WHERE email = $1 LIMIT 1",
    ufd.User.Email,
  )

  // <email>
  if user.Email == "" {
    r.JSON(200, map[string]string{ "error": "User not found." })
    return
  }

  // <password>
  bcrypt_check_err := bcrypt.CompareHashAndPassword(
    []byte(user.EncryptedPassword),
    []byte(ufd.User.Password),
  )

  if bcrypt_check_err != nil {
    r.JSON(200, map[string]string{ "error": "Invalid password." })
    return
  }

  // <success>
  token := GenerateToken(&user)
  user_public := UserPublic{ Token: token }

  r.JSON(200, map[string]UserPublic{ "user": user_public })
}



//
//  -> VERIFY TOKEN
//
func Users__VerifyToken(req *http.Request, r render.Render) {
  qs := req.URL.Query()
  is_valid := VerifyToken(qs.Get("token"))

  r.JSON(200, map[string]bool{ "is_valid": is_valid })
}
