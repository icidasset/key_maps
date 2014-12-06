package main

import (
  "github.com/go-martini/martini"
  "github.com/icidasset/key-maps/api"
  "github.com/icidasset/key-maps/db"
  "github.com/martini-contrib/binding"
  "github.com/martini-contrib/render"
  "html/template"
  "io/ioutil"
  "net/http"
)


//
//  [Root]
//  -> HTML files (for js application)
//
func rootHandler(w http.ResponseWriter) {
  files, _ := ioutil.ReadDir("./views/templates/")
  filepaths := make([]string, len(files))

  for _, f := range files {
    p := "views/templates/" + f.Name()
    filepaths = append(filepaths, p)
  }

  tmpl, _ := template.ParseFiles(
    "views/layout.html",
    "views/index.html",
  )

  tmpl.ExecuteTemplate(w, "layout", filepaths)
}


//
//  [Main]
//
func main() {
  r := martini.Classic()
  r.Use(render.Renderer())

  // prepare database
  if err := db.Open(); err != nil {
    panic(err)
  }

  defer db.Close()

  // routes
  r.Get("/api/maps", api.Maps__Index)
  r.Get("/api/maps/:id", api.Maps__Show)
  r.Post("/api/maps", binding.Bind(api.MapFormData{}), api.Maps__Create)

  r.Get("/", rootHandler)

  // setup server
  r.Run()
}
