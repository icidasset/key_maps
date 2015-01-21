package main

import (
  "github.com/gocraft/web"
  "github.com/icidasset/key-maps/api"
  "net/http"
  "net/http/httptest"
  "strings"
  "testing"
)

/*

  TODO:

  - gzip
  - √ bindings
  - √ must_be_authenticated middleware
  - √ root_handler -> render html
  - api handlers -> render json
  - CORS (public routes only)

*/

func TestRootHandler(t *testing.T) {
  router := web.New(api.Context{})
  CreateRootRoute(router)

  rr, req := newTestRequest("GET", "/")
  router.ServeHTTP(rr, req)

  // should load route and contain ember templates
  substr := `<script type="text/x-handlebars" data-template-name="index">`

  if rr.Code != 200 {
    t.Error("Could not load the root route.")
  } else if strings.Contains(responseBodyToString(rr), substr) == false {
    t.Error("Could not find the ember templates when making request to the root route.")
  }
}



//
//  Testing helpers
//  -> most are taken from the gocraft/web project
//
func newTestRequest(method, path string) (*httptest.ResponseRecorder, *http.Request) {
  request, _ := http.NewRequest(method, path, nil)
  recorder := httptest.NewRecorder()

  return recorder, request
}


func responseBodyToString(rr *httptest.ResponseRecorder) string {
  return strings.TrimSpace(string(rr.Body.Bytes()))
}
