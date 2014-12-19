package main

import (
  "github.com/go-martini/martini"
  "github.com/icidasset/key-maps/api"
  "github.com/icidasset/key-maps/db"
  "github.com/martini-contrib/binding"
  "github.com/martini-contrib/render"
  "io/ioutil"
  "net/http"
  "strings"
  "text/template"
)


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


func MustBeAuthenticatedMiddleware(c martini.Context, w http.ResponseWriter, r *http.Request) {
  auth_header := r.Header.Get("Authorization")
  t := strings.Split(auth_header, "Bearer ")[1]
  token := api.ParseToken(t)

  if !token.Valid {
    http.Error(w, "Forbidden", http.StatusUnauthorized)
  } else {
    id := int(token.Claims["user_id"].(float64))
    c.Map(api.User{ Id: id })
  }
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

  // routes
  r.Group("/api/users", func(r martini.Router) {
    r.Get("/verify-token", api.Users__VerifyToken)

    r.Post(
      "",
      binding.Bind(api.UserAuthFormData{}),
      api.Users__Create,
    )

    r.Post(
      "/authenticate",
      binding.Bind(api.UserAuthFormData{}),
      api.Users__Authenticate,
    )
  })

  r.Group("/api/maps", func(r martini.Router) {
    r.Get("", api.Maps__Index)
    r.Get("/:id", api.Maps__Show)

    r.Post(
      "",
      binding.Bind(api.MapFormData{}),
      api.Maps__Create,
    )
  }, MustBeAuthenticatedMiddleware)

  r.Get("/", rootHandler)

  // setup server
  r.Run()
}
