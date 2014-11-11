package main

import (
	"github.com/icidasset/key-maps/api"
	_ "github.com/lib/pq"
	"github.com/pilu/traffic"
	"html/template"
)


//
//  [Root]
//  -> HTML files (for js application)
//
func rootHandler(w traffic.ResponseWriter, r *traffic.Request) {
	tmpl, _ := template.ParseFiles(
		"views/layout.html",
		"views/index.html",
	)

	tmpl.ExecuteTemplate(w, "layout", nil)
}


//
//  [Errors]
//
func notFoundHandler(w traffic.ResponseWriter, r *traffic.Request) {
	w.WriteText("404")
}


func errorHandler(w traffic.ResponseWriter, r *traffic.Request, err interface {}) {
	w.WriteText("500")
}


//
//  [Main]
//
func main() {
	r := traffic.New()

	// routes
	r.Get("/api/maps/?", api.GetMaps)
	r.Get("/", rootHandler)

	// errors
	r.NotFoundHandler = notFoundHandler
	r.ErrorHandler = errorHandler

	// setup server
	r.Run()
}
