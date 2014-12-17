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
  Id int                      `json:"id"`
  Email string                `json:"email"`
  EncryptedPassword string    `json:"-"`
  CreatedAt time.Time         `json:"created_at" db:"created_at"`
  UpdatedAt time.Time         `json:"updated_at" db:"updated_at"`
}


type UserNew struct {
  Email string                `form:"email" binding:"required"`
  Password string             `form:"password" binding:"required"`
  PasswordConfirmation string `form:"password_confirmation" binding:"required"`
}


type UserNewFormData struct {
  User UserNew `form:"user" binding:"required"`
}


type UserAuth struct {
  Email string                `form:"email" binding:"required"`
  Password string             `form:"password" binding:"required"`
}


type UserAuthFormData struct {
  User UserAuth `form:"user" binding:"required"`
}



//
//  Routes
//
func Users__Create(ufd UserNewFormData, r render.Render) {
  query := "INSERT INTO users (email, encrypted_password, created_at, updated_at) VALUES (:email, :encrypted_password, :created_at, :updated_at)"

  // make new user
  encryped_password, _ := bcrypt.GenerateFromPassword(
     []byte(ufd.User.Password),
     bcrypt.DefaultCost,
  )
  now := time.Now()

  new_user := User{Email: ufd.User.Email, EncryptedPassword: string(encryped_password), CreatedAt: now, UpdatedAt: now}

  // execute query
  _, err := db.Inst().NamedExec(query, new_user)

  // if error
  if err != nil {
    r.JSON(500, err.Error())

  // render map as json
  } else {
    u := User{}
    db.Inst().Get(&u, "SELECT * FROM users WHERE email = $1 LIMIT 1", new_user.Email)
    r.JSON(200, map[string]User{ "user": u })

  }
}


func Users__Authenticate(ufd UserAuthFormData, r render.Render) {
  user := User{}
  db.Inst().Get(&user, "SELECT * FROM users WHERE email = $1 LIMIT 1", ufd.User.Email)

  if user.Email == "" {
    // STOP, user doesn't exist
  }

  bcrypt_check_err := bcrypt.CompareHashAndPassword(
    []byte(user.EncryptedPassword),
    []byte(ufd.User.Password),
  )

  if bcrypt_check_err != nil {
    // STOP, invalid password
  }

  token := jwt.New(jwt.GetSigningMethod("HS256"))
  token.Claims["user_id"] = user.Id
  token.Claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
  token_string, _ := token.SignedString([]byte("TODO - SECRET KEY"))

  r.JSON(200, map[string]string{ "token" : token_string })
}
