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
  "os"
)


var SecretKey string


type TemplateData struct {
  EmberTemplates string
}


func ScanTemplatesDir(path string) string {
  files, _ := ioutil.ReadDir(path)
  templates := make([]string, 0)

  for _, f := range files {
    name := f.Name()

    if f.IsDir() {
      t := ScanTemplatesDir(path + name + "/")
      templates = append(templates, t)

    } else if strings.HasSuffix(name, ".html") {
      c, _ := ioutil.ReadFile(path + name)
      templates = append(templates, string(c))
    }
  }

  return strings.Join(templates, "\n")
}



//
//  [Root]
//  -> HTML files (for js application)
//
func rootHandler(w http.ResponseWriter) {
  tmpl, _ := template.ParseFiles(
    "views/layout.html",
    "views/index.html",
  )

  ember_templates := ScanTemplatesDir("views/ember_templates/")
  tmpl_data := TemplateData{ EmberTemplates: ember_templates }
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

  // environment variables
  SecretKey = os.Getenv("SECRET_KEY")

  // routes
  r.Post("/api/users", binding.Bind(api.UserAuthFormData{}), api.Users__Create)
  r.Post("/api/users/authenticate", binding.Bind(api.UserAuthFormData{}), api.Users__Authenticate)

  r.Get("/api/maps", api.Maps__Index)
  r.Get("/api/maps/:id", api.Maps__Show)
  r.Post("/api/maps", binding.Bind(api.MapFormData{}), api.Maps__Create)

  r.Get("/", rootHandler)

  // setup server
  r.Run()
}
