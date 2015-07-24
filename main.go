package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/icidasset/key-maps-api/api"
	"github.com/icidasset/key-maps-api/db"
	"github.com/icidasset/key-maps-api/middleware"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

func main() {
	env := os.Getenv("ENV")
	host := os.Getenv("HOST")

	// flags
	port := flag.String("port", "3000", "Server port address")

	flag.Parse()

	// new router
	router := echo.New()

	router.Use(mw.Recover())
	router.Use(mw.Gzip())
	router.Use(middleware.Cors)

	// extra middleware
	if env == "" || env == "development" {
		router.Use(mw.Logger())
	}

	// prepare database
	if err := db.Open(); err != nil {
		panic(err)
	}

	defer db.Close()

	// routes
	CreateUserRoutes(router)
	CreateMapRoutes(router)
	CreateMapItemRoutes(router)
	CreatePublicRoutes(router)

	// run
	if host == "" {
		host = "0.0.0.0"
	}

	http.ListenAndServe(host+":"+*port, router)
}

//
//  Routes — Users
//
func CreateUserRoutes(router *echo.Echo) {
	g := router.Group("/users")

	g.Get("/verify-token", api.Users__VerifyToken)
	g.Post("", api.Users__Create)
	g.Post("/authenticate", api.Users__Authenticate)

	// option requests
	g.Options("", Options)
	g.Options("/verify-token", Options)
	g.Options("/authenticate", Options)
}

//
//  Routes — Maps
//
func CreateMapRoutes(router *echo.Echo) {
	g := router.Group("/maps")
	g.Use(middleware.MustBeAuthenticated)

	g.Get("", api.Maps__Index)
	g.Get("/:id", api.Maps__Show)
	g.Delete("/:id", api.Maps__Destroy)
	g.Post("", api.Maps__Create)
	g.Patch("/:id", api.Maps__Update)
	g.Put("/:id", api.Maps__Update)

	// option requests
	g.Options("", Options)
	g.Options("/:id", Options)
}

//
//  Routes — Map Items
//
func CreateMapItemRoutes(router *echo.Echo) {
	g := router.Group("/map_items")
	g.Use(middleware.MustBeAuthenticated)

	g.Get("/:id", api.MapItems__Show)
	g.Delete("/:id", api.MapItems__Destroy)
	g.Post("", api.MapItems__Create)
	g.Patch("/:id", api.MapItems__Update)
	g.Put("/:id", api.MapItems__Update)

	// option requests
	g.Options("", Options)
	g.Options("/:id", Options)
}

//
//  Routes — Public
//
func CreatePublicRoutes(router *echo.Echo) {
	g := router.Group("/public")

	g.Get("/:hash", api.Public__Show)

	// option requests
	g.Options("", Options)
}

//
//  Options
//
func Options(c *echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, PATCH, DELETE")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")

	return c.NoContent(200)
}
