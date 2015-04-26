package api

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/icidasset/key-maps-api/db"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id                int
	Email             string
	EncryptedPassword string    `db:"encrypted_password"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

type UserAuth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserAuthFormData struct {
	User UserAuth `json:"user"`
}

type UserPublic struct {
	Token string `json:"token"`
}

type UserHandlers struct {
	Test string
}

//
//  {post} CREATE
//
func Users__Create(c *echo.Context) {
	query := "INSERT INTO users (email, encrypted_password, created_at, updated_at) VALUES (:email, :encrypted_password, :created_at, :updated_at) RETURNING id"

	// parse json from request body
	uafd := UserAuthFormData{}
	json_decoder := json.NewDecoder(c.Request.Body)
	json_decoder.Decode(&uafd)

	// make new user
	encryped_password, _ := bcrypt.GenerateFromPassword(
		[]byte(uafd.User.Password),
		bcrypt.DefaultCost,
	)

	now := time.Now()

	new_user := User{
		Email:             strings.ToLower(uafd.User.Email),
		EncryptedPassword: string(encryped_password),
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	// execute query
	rows, err := db.Inst().NamedQuery(query, new_user)

	// return if error
	if err != nil {
		c.JSON(500, FormatError(err))
		return
	}

	// scan rows
	for rows.Next() {
		err = rows.StructScan(&new_user)
	}

	// generate token for user
	token := GenerateToken(&new_user)
	user_public := UserPublic{Token: token}

	// render
	if err != nil {
		c.JSON(500, FormatError(err))
	} else {
		c.JSON(201, map[string]UserPublic{"user": user_public})
	}
}

//
//  {post} AUTHENTICATE
//
func Users__Authenticate(c *echo.Context) {
	user := User{}

	// parse json from request body
	uafd := UserAuthFormData{}
	json_decoder := json.NewDecoder(c.Request.Body)
	json_decoder.Decode(&uafd)

	// query
	db.Inst().Get(
		&user,
		"SELECT * FROM users WHERE email = $1",
		strings.ToLower(uafd.User.Email),
	)

	// <email>
	if user.Email == "" {
		c.JSON(200, map[string]string{"error": "User not found."})
		return
	}

	// <password>
	bcrypt_check_err := bcrypt.CompareHashAndPassword(
		[]byte(user.EncryptedPassword),
		[]byte(uafd.User.Password),
	)

	if bcrypt_check_err != nil {
		c.JSON(200, map[string]string{"error": "Invalid password."})
		return
	}

	// <success>
	token := GenerateToken(&user)
	user_public := UserPublic{Token: token}

	c.JSON(200, map[string]UserPublic{"user": user_public})
}

//
//  {get} VERIFY TOKEN
//
func Users__VerifyToken(c *echo.Context) {
	qs := c.Request.URL.Query()
	token, err := ParseToken(qs.Get("token"))
	is_valid := false

	if err == nil && token.Valid {
		is_valid = true
	}

	// invalid token
	if !is_valid {
		c.JSON(200, map[string]bool{"is_valid": false})

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
			c.JSON(200, map[string]bool{"is_valid": true})

			// user does not exist
		} else {
			c.JSON(200, map[string]bool{"is_valid": false})

		}

	}
}
