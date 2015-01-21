package main

import (
  "github.com/gocraft/web"
  "github.com/icidasset/key-maps/api"
  "github.com/icidasset/key-maps/db"
  "github.com/opennota/json-binding"
  "io/ioutil"
  "net/http"
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
  router := web.New(api.BaseContext{})
  router.Middleware(web.StaticMiddleware("public"))

  // prepare database
  if err := db.Open(); err != nil {
    panic(err)
  }

  defer db.Close()

  // routes
  CreateRootRoute(router)
  CreateUserRoutes(router)
  // CreateMapRoutes(router)
  // CreateMapItemRoutes(router)
  // CreatePublicRoutes(router)

  // run
  http.ListenAndServe("localhost:3000", router)
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

    Middleware(binding.Bind(api.UserAuthFormData{})).

    Post("", (*api.Context).Users__Create).
    Post("/authenticate", (*api.Context).Users__Authenticate)
}


//
//  Routes — Maps
//
// func CreateMapRoutes(router *web.Router) {
//   api_maps_router = router.Subrouter(ApiMapsContext{}, "/api/maps").
//     Middleware((*ApiMapItemsContext).MustBeAuthenticated).
//
//     Get("", (*ApiMapsContext).api.Maps__Index).
//     Get("/:id", (*ApiMapsContext).api.Maps__Show).
//     Delete("/:id", (*ApiMapsContext).api.Maps__Destroy).
//
//     Middleware(binding.Bind(api.MapFormData{})).
//
//     Post("", (*ApiMapsContext).api.Maps__Create).
//     Put("/:id", (*ApiMapsContext).api.Maps__Update)
// }


//
//  Routes — Map Items
//
// func CreateMapItemRoutes(router *web.Router) {
//   api_map_items_router = router.Subrouter(ApiMapItemsContext{}, "/api/map_items").
//     Middleware((*ApiMapItemsContext).MustBeAuthenticated).
//
//     Get("/:id", (*ApiMapItemsContext).api.MapItems__Show).
//     Delete("/:id", (*ApiMapItemsContext).api.MapItems__Destroy).
//
//     Middleware(binding.Bind(api.MapItemFormData{})).
//
//     Post("", (*ApiMapItemsContext).api.MapItems__Create).
//     Put("/:id", (*ApiMapItemsContext).api.MapItems__Update)
// }


//
//  Routes — Public
//
// func CreatePublicRoutes(router *web.Router) {
//   api_public_router = router.Subrouter(ApiPublicContext{}, "/api/public").
//     Get("/:hash", (*ApiPublicContext).api.Public__Show)
// }
