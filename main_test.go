package main

import (
  "bytes"
  "encoding/json"
  "fmt"
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



//
//  Root
//
func TestRootHandler(t *testing.T) {
  router := web.New(api.Context{})
  CreateRootRoute(router)

  rec, req := newTestRequest("GET", "/")
  router.ServeHTTP(rec, req)

  // should load route and contain ember templates
  substr := `<script type="text/x-handlebars" data-template-name="index">`

  if rec.Code != 200 {
    t.Error("Could not load the root route.")
  } else if strings.Contains(responseBodyToString(rec), substr) == false {
    t.Error("Could not find the ember templates when making request to the root route.")
  }
}



//
//  API - Users
//
func TestApiUsers(t *testing.T) {
  router := web.New(api.Context{})
  CreateUserRoutes(router)

  // create user
  user := api.UserAuth{ Email: "test@gmail.com", Password: "password" }
  user_form_data := api.UserAuthFormData{ User: user }
  j, _ := json.Marshal(user_form_data)

  req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(j))
  req.Header.Set("Content-Type", "application/json")
  rec := httptest.NewRecorder()
  router.ServeHTTP(rec, req)

  // result := map[string]api.UserPublic{}
  // json.Unmarshal(rec.Body.Bytes(), result)

  if rec.Code != 201 {
    t.Error("Did not create user correctly.")
  }

  // else if result["user"].Token == "" {
  //   t.Error("Did not return the user's token on user create.")
  // }
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
