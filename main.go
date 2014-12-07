package main

import (
  "github.com/go-martini/martini"
  "github.com/icidasset/key-maps/api"
  "github.com/icidasset/key-maps/db"
  "github.com/martini-contrib/binding"
  "github.com/martini-contrib/render"
  "strings"
  "text/template"
  "io/ioutil"
  "net/http"
)

type TemplateData struct {
  EmberTemplates string
}


//
//  [Root]
//  -> HTML files (for js application)
//
func rootHandler(w http.ResponseWriter) {
  et_base_path := "views/ember_templates/"
  et_files, _ := ioutil.ReadDir(et_base_path)
  et_file_contents := make([]string, 0)

  for _, f := range et_files {
    name := f.Name()
    if strings.HasSuffix(name, ".html") {
      c, _ := ioutil.ReadFile(et_base_path + name)
      et_file_contents = append(et_file_contents, string(c))
    }
  }

  et := strings.Join(et_file_contents, "\n")

  tmpl, _ := template.ParseFiles(
    "views/layout.html",
    "views/index.html",
  )

  tmpl_data := TemplateData{ EmberTemplates: et }
  tmpl.ExecuteTemplate(w, "layout", tmpl_data)
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
