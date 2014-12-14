package api

import (
  _ "github.com/go-martini/martini"
  _ "github.com/lib/pq"
  "github.com/icidasset/key-maps/db"
  "github.com/martini-contrib/render"
  "github.com/satori/go.uuid"
  "golang.org/x/crypto/bcrypt"
  "time"
)

type User struct {
  Id int                      `json:"id"`
  Email string                `json:"email"`
  EncryptedPassword string    `json:"-"`
  AccessToken string          `json:"access_token"`
  CreatedAt time.Time         `json:"created_at" db:"created_at"`
  UpdatedAt time.Time         `json:"updated_at" db:"updated_at"`
}


type UserNew struct {
  Email string                `form:"email" binding:"required"`
  Password string             `form:"password" binding:"required"`
  PasswordConfirmation string `form:"password_confirmation" binding:"required"`
}


type UserFormData struct {
  User UserNew `form:"user" binding:"required"`
}


func Users__Create(ufd UserFormData, r render.Render) {
  query := "INSERT INTO users (email, encrypted_password, created_at, updated_at) VALUES (:email, :encrypted_password, :created_at, :updated_at)"

  // make new user
  encryped_password, _ := bcrypt.GenerateFromPassword(
     []byte(ufd.User.Password),
     bcrypt.DefaultCost,
  )
  access_token := uuid.NewV4().String();
  now := time.Now()

  new_user := User{Email: ufd.User.Email, EncryptedPassword: string(encryped_password), AccessToken: access_token, CreatedAt: now, UpdatedAt: now}

  // execute query
  _, err := db.Inst().NamedExec(query, new_user)

  // if error
  if err != nil {
    r.JSON(500, err.Error())

  // render map as json
  } else {
    u := User{}
    db.Inst().Get(&u, "SELECT * FROM users WHERE email = $1", new_user.Email)
    r.JSON(200, map[string]User{ "user": u })

  }
}
