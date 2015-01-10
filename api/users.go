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
//  {post} CREATE
//
func Users__Create(ufd UserAuthFormData, r render.Render) {
  query := "INSERT INTO users (email, encrypted_password, created_at, updated_at) VALUES (:email, :encrypted_password, :created_at, :updated_at) RETURNING id"

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
  rows, err := db.Inst().NamedQuery(query, new_user)

  // return if error
  if err != nil {
    r.JSON(500, FormatError(err));
    return
  }

  // scan rows
  for rows.Next() {
    err = rows.StructScan(&new_user)
  }

  // generate token for user
  token := GenerateToken(&new_user)
  user_public := UserPublic{ Token: token }

  // render
  if err != nil {
    r.JSON(500, FormatError(err));
  } else {
    r.JSON(201, map[string]UserPublic{ "user": user_public })
  }
}



//
//  {post} AUTHENTICATE
//
func Users__Authenticate(ufd UserAuthFormData, r render.Render) {
  user := User{}

  db.Inst().Get(
    &user,
    "SELECT * FROM users WHERE email = $1",
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
//  {get} VERIFY TOKEN
//
func Users__VerifyToken(req *http.Request, r render.Render) {
  qs := req.URL.Query()
  token := ParseToken(qs.Get("token"))
  is_valid := token.Valid

  // invalid token
  if !is_valid {
    r.JSON(200, map[string]bool{ "is_valid": false })

  // valid token, but check if the user exists
  } else {
    user := User{}

    db.Inst().Get(
      &user,
      "SELECT id FROM users WHERE id = $1",
      int(token.Claims["user_id"].(float64)),
    )

    // user exists
    if user.Id != 0 {
      r.JSON(200, map[string]bool{ "is_valid": true })

    // user does not exist
    } else {
      r.JSON(200, map[string]bool{ "is_valid": false })

    }

  }
}
