package main

import (
  "github.com/icidasset/key-maps/api"
  "github.com/icidasset/key-maps/db"
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

  // prepare database
  if err := db.Open(); err != nil {
    panic(err)
  }

  defer db.Close()

  // routes
  r.Get("/api/maps", api.GetMaps)
  r.Get("/api/maps/:id", api.GetMap)
  r.Get("/", rootHandler)

  // errors
  r.NotFoundHandler = notFoundHandler
  r.ErrorHandler = errorHandler

  // setup server
  r.Run()
}
