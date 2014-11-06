package main

import (
	"fmt"
	"github.com/gocraft/web"
	"html/template"
	"net/http"
	"path"
)


//
//  [Context]
//
type Context struct {
}


//
//  [Root]
//
func Root(w web.ResponseWriter, r *web.Request) {
	file := "index.html" // r.URL.Path

	lp := path.Join("templates", "layout.html")
	fp := path.Join("templates", file)

	tmpl, _ := template.ParseFiles(lp, fp)
	tmpl.ExecuteTemplate(w, "layout", nil)
}


//
//  [Errors]
//
func (c *Context) NotFound(w web.ResponseWriter, r *web.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Not Found")
}


//
//  [Main]
//
func main() {
	router := web.New(Context {}).
		NotFound((*Context).NotFound).
		Middleware(web.ShowErrorsMiddleware).
		Middleware(web.StaticMiddleware("public")).
		Get("/", Root)

	http.ListenAndServe("localhost:3000", router)
}
