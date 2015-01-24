package main

import (
  "flag"
  "github.com/gocraft/web"
  "github.com/icidasset/key-maps/api"
  "github.com/icidasset/key-maps/db"
  "io/ioutil"
  "net/http"
  "os"
  "strings"
  "text/template"
)


//
//  [Root]
//  -> HTML files (for js application)
//
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


func rootHandler(rw web.ResponseWriter, req *web.Request) {
  tmpl, _ := template.ParseFiles(
    "views/layout.html",
    "views/index.html",
  )

  ember_templates := ScanTemplatesDir("views/ember_templates/")
  tmpl_data := TemplateData{ EmberTemplates: ember_templates }
  tmpl.ExecuteTemplate(rw, "layout", tmpl_data)
}


//
//  [Main]
//
func main() {
  env := os.Getenv("ENV")

  // flags
  port := flag.String("port", "3000", "Server port address")

  flag.Parse()

  // new router
  router := web.New(api.BaseContext{})
  router.Middleware((*api.BaseContext).Gzip)
  router.Middleware(web.StaticMiddleware("public"))

  // extra middleware
  if env == "" || env == "development" {
    router.Middleware(web.LoggerMiddleware)
    router.Middleware(web.ShowErrorsMiddleware)
  }

  // prepare database
  if err := db.Open(); err != nil {
    panic(err)
  }

  defer db.Close()

  // routes
  CreateRootRoute(router)
  CreateUserRoutes(router)
  CreateMapRoutes(router)
  CreateMapItemRoutes(router)
  CreatePublicRoutes(router)

  // run
  http.ListenAndServe("localhost:" + *port, router)
}


//
//  Routes — Root
//
func CreateRootRoute(router *web.Router) {
  router.Get("/", rootHandler)
}


//
//  Routes — Users
//
func CreateUserRoutes(router *web.Router) {
  router.Subrouter(api.Context{}, "/api/users").
    Get("/verify-token", (*api.Context).Users__VerifyToken).

    Post("/", (*api.Context).Users__Create).
    Post("/authenticate", (*api.Context).Users__Authenticate)
}


//
//  Routes — Maps
//
func CreateMapRoutes(router *web.Router) {
  router.Subrouter(api.Context{}, "/api/maps").
    Middleware((*api.Context).MustBeAuthenticated).

    Get("/", (*api.Context).Maps__Index).
    Get("/:id", (*api.Context).Maps__Show).
    Delete("/:id", (*api.Context).Maps__Destroy).

    Post("/", (*api.Context).Maps__Create).
    Put("/:id", (*api.Context).Maps__Update)
}


//
//  Routes — Map Items
//
func CreateMapItemRoutes(router *web.Router) {
  router.Subrouter(api.Context{}, "/api/map_items").
    Middleware((*api.Context).MustBeAuthenticated).

    Get("/:id", (*api.Context).MapItems__Show).
    Delete("/:id", (*api.Context).MapItems__Destroy).

    Post("", (*api.Context).MapItems__Create).
    Put("/:id", (*api.Context).MapItems__Update)
}


//
//  Routes — Public
//
func CreatePublicRoutes(router *web.Router) {
  router.Subrouter(api.Context{}, "/api/public").
    Middleware((*api.Context).CORS).

    Get("/:hash", (*api.Context).Public__Show)
}
