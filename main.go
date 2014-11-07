package main

import (
	"database/sql"
	"github.com/icidasset/key-maps/api"
	_ "github.com/lib/pq"
	"github.com/pilu/traffic"
	"html/template"
	"path"
)


//
//  [Root]
//  -> HTML files (for js application)
//
func rootHandler(w traffic.ResponseWriter, r *traffic.Request) {
	file := "index.html"

	layoutPath := path.Join("templates", "layout.html")
	filePath := path.Join("templates", file)

	tmpl, _ := template.ParseFiles(layoutPath, filePath)
	tmpl.ExecuteTemplate(w, "layout", nil)
}


//
//  [Errors]
//
func notFoundHandler(w traffic.ResponseWriter, r *traffic.Request) {
	w.Render("404")
}


func errorHandler(w traffic.ResponseWriter, r *traffic.Request, err interface {}) {
	w.Render("500")
}


//
//  [Main]
//
func main() {
	db, _ := sql.Open("postgres", "user=icidasset dbname=keymaps_development")
	defer db.Close()

	// router
	// -> serves static files from "public" directory
	r := traffic.New()

	// routes
	r.Get("/api/lists/?", api.GetLists)
	r.Get("/", rootHandler)

	// errors
	r.NotFoundHandler = notFoundHandler
	r.ErrorHandler = errorHandler

	// setup server
	r.Run()
}
