package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocraft/web"
	"html/template"
	"net/http"
	"path"
)


//
//  [Contexts]
//
type Context struct {
}


//
//  [API]
//  -> https://gobyexample.com/json
//  -> http://www.alexedwards.net/blog/golang-response-snippets
//  -> http://godoc.org/github.com/lib/pq
//
func ApiGetLists(w web.ResponseWriter, r *web.Request) {
	m := make(map[string]string)
	m["an"] = "example"
	m["a"] = "test"

	j, _ := json.Marshal(m)

	w.Header().Set("Content-Type", "application/json")
  w.Write(j)
}


//
//  [Root]
//  -> HTML files (for js application)
//
func Root(w web.ResponseWriter, r *web.Request) {
	file := "index.html" // r.URL.Path

	layoutPath := path.Join("templates", "layout.html")
	filePath := path.Join("templates", file)

	tmpl, _ := template.ParseFiles(layoutPath, filePath)
	tmpl.ExecuteTemplate(w, "layout", nil)
}


//
//  [Errors]
//
func NotFound(w web.ResponseWriter, r *web.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Not Found")
}


//
//  [Main]
//  -> https://developer.uservoice.com/blog/2013/12/12/why-i-wrote-gocraft-web/
//
func main() {
	r := web.New(Context {})
	r.NotFound(NotFound)
	r.Middleware(web.ShowErrorsMiddleware)
	r.Middleware(web.StaticMiddleware("public"))

	// routes
	r.Get("/api/lists/", ApiGetLists)
	r.Get("/", Root)

	// setup server
	http.ListenAndServe("localhost:3000", r)
}
