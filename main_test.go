package main

import (
  "bytes"
  "encoding/json"
  "github.com/gocraft/web"
  "github.com/icidasset/key-maps/api"
  "github.com/icidasset/key-maps/db"
  . "gopkg.in/check.v1"
  "io"
  "net/http"
  "net/http/httptest"
  "os"
  "os/exec"
  "strconv"
  "strings"
  "testing"

  "fmt"
)

/*

  TODO:

  - gzip
  - √ bindings
  - √ must_be_authenticated middleware
  - √ root_handler -> render html
  - √ api handlers -> render json
  - CORS (public routes only)

*/

// setup gocheck
func Test(t *testing.T) { TestingT(t) }

type MySuite struct {
  router *web.Router

  // users
  userAuthToken string

  // maps
  mapId int
}

var _ = Suite(&MySuite{})


func (s *MySuite) SetUpSuite(c *C) {
  os.Setenv("ENV", "test")
  os.Setenv("SECRET_KEY", "doesntmatter")

  // database
  create_db := exec.Command("createdb", "keymaps_test")
  create_db.Run()

  if err := db.Open(); err != nil {
    panic(err)
  }

  // router
  s.router = web.New(api.Context{})
  CreateRootRoute(s.router)
  CreateUserRoutes(s.router)
  CreateMapRoutes(s.router)
}


func (s *MySuite) TearDownSuite(c *C) {
  db.Close()

  drop_db := exec.Command("dropdb", "keymaps_test")
  drop_db.Run()
}


func (s *MySuite) TestAll(c *C) {
  (s).testRootHandler(c)
  (s).testApiUsers(c)
  (s).testApiMaps(c)
}



//
//  Root
//
func (s *MySuite) testRootHandler(c *C) {
  req, rec := newTestRequest("GET", "/", emptyBuffer())
  s.router.ServeHTTP(rec, req)

  // should load route and contain ember templates
  substr := `<script type="text/x-handlebars" data-template-name="index">`

  if rec.Code != 200 {
    c.Error("Could not load the root route.")
  } else if strings.Contains(responseBodyToString(rec), substr) == false {
    c.Error("Could not find the ember templates when making request to the root route.")
  }
}



//
//  API - Users
//
func (s *MySuite) testApiUsers(c *C) {
  (s).testApiUsers__Create(c)
  (s).testApiUsers__Authenticate(c)
  (s).testApiUsers__VerifyToken(c)
}


/*

  POST '/api/users'

  @data --- request-body (application/json)
    { "user": { "email": "test@gmail.com", "password": "password" } }

  @return --- response-body (application/json)
    { "token": "some-generated-token" }

*/
func (s *MySuite) testApiUsers__Create(c *C) {
  user := api.UserAuth{ Email: "test@gmail.com", Password: "password" }
  user_form_data := api.UserAuthFormData{ User: user }
  j, _ := json.Marshal(user_form_data)

  // make request
  req, rec := newTestRequest("POST", "/api/users", bytes.NewBuffer(j))
  req.Header.Set("Content-Type", "application/json")
  s.router.ServeHTTP(rec, req)

  // parse json from response
  result := map[string]api.UserPublic{}
  json.Unmarshal(rec.Body.Bytes(), &result)

  // validate
  if rec.Code != 201 {
    c.Error("Did not create user correctly.")
  } else if result["user"].Token == "" {
    c.Error("Did not return the user's token on user create.")
  }
}


/*

  POST '/api/users/authenticate'

  @data --- request-body (application/json)
    { "user": { "email": "test@gmail.com", "password": "password" } }

  @return --- response-body (application/json)
    { "token": "some-generated-token" }

*/
func (s *MySuite) testApiUsers__Authenticate(c *C) {
  user := api.UserAuth{ Email: "test@gmail.com", Password: "password" }
  user_form_data := api.UserAuthFormData{ User: user }
  j, _ := json.Marshal(user_form_data)

  // make request
  req, rec := newTestRequest("POST", "/api/users/authenticate", bytes.NewBuffer(j))
  req.Header.Set("Content-Type", "application/json")
  s.router.ServeHTTP(rec, req)

  // parse json from response
  result := map[string]api.UserPublic{}
  err := json.Unmarshal(rec.Body.Bytes(), &result)

  // validate
  if err != nil {
    c.Error("Did not authenticate correctly.")
  } else if result["user"].Token == "" {
    c.Error("Did not return the user's token when trying to authenticate.")
  } else {
    s.userAuthToken = result["user"].Token
  }
}


/*

  GET '/api/users/verify-token?token=TOKEN_FROM_AUTHENTICATE_REQUEST'

  @data --- query-string-param (string)
    'token'

  @return --- response-body (application/json)
    { "is_valid": boolean }

*/
func (s *MySuite) testApiUsers__VerifyToken(c *C) {
  url := "/api/users/verify-token?token=" + s.userAuthToken

  req, rec := newTestRequest("GET", url, emptyBuffer())
  s.router.ServeHTTP(rec, req)

  // parse json from response
  result := map[string]bool{}
  err := json.Unmarshal(rec.Body.Bytes(), &result)

  // validate
  if err != nil {
    c.Error("Something went wrong when trying to verify token.")
  } else if result["is_valid"] != true {
    c.Error("Could not verify token.")
  }
}



//
//  API - Maps
//
func (s *MySuite) testApiMaps(c *C) {
  (s).testApiMaps__Create(c)
}


/*

  POST '/api/maps'

  @data --- request-body (application/json)
    { "map": { "name": "Quotes", "structure": "[]" } }

  @return --- response-body (application/json)
    { "map": { id: "...", "name": "Quotes", "slug": "quotes", "structure": "[]", created_at: "..." } }

*/
func (s *MySuite) testApiMaps__Create(c *C) {
  m := api.Map{ Name: "Quotes", Structure: "[]" }
  m_form_data := api.MapFormData{ Map: m }
  j, _ := json.Marshal(m_form_data)

  fmt.Print("\n\n\n")
  fmt.Printf("%v", string(j));
  fmt.Print("\n\n\n")

  // make request
  req, rec := newTestRequest("POST", "/api/maps", bytes.NewBuffer(j))
  req.Header.Set("Content-Type", "application/json")
  setAuthorizationHeader(req, s)
  s.router.ServeHTTP(rec, req)

  // parse json from response
  result := map[string]api.Map{}
  json.Unmarshal(rec.Body.Bytes(), &result)

  fmt.Print("\n\n\n")
  fmt.Printf("%#v", result["map"])
  fmt.Print("\n\n\n")

  // validate
  if rec.Code != 201 {
    c.Error("Did not create map correctly.")
  } else if result["map"].Id == 0 {
    c.Error("Did not return a map.")
  } else if result["map"].Slug != "quotes" {
    c.Error("Incorrect generation of the slug.")
  } else {
    s.mapId = result["map"].Id
  }
}


/*

  PUT '/api/maps/:id'

  @data --- request-body (application/json)
    { "map": { "name": "Quotes", "structure": "[]" } }

  @return --- response-body (application/json)
    { "map": { id: "...", "name": "Quotes", "slug": "quotes", "structure": "[]", created_at: "..." } }

*/
func (s *MySuite) testApiMaps__Update(c *C) {
  new_structure := `[{ "key": "test", "type": "string" }]`
  new_sort_by := `test`

  // make json
  m := api.Map{ Name: "Quotes", Structure: new_structure, SortBy: new_sort_by }
  m_form_data := api.MapFormData{ Map: m }
  j, _ := json.Marshal(m_form_data)

  // make request
  req, rec := newTestRequest("POST", "/api/maps/" + strconv.Itoa(s.mapId), bytes.NewBuffer(j))
  req.Header.Set("Content-Type", "application/json")
  setAuthorizationHeader(req, s)
  s.router.ServeHTTP(rec, req)

  // parse json from response
  result := map[string]api.Map{}
  json.Unmarshal(rec.Body.Bytes(), &result)

  // validate
  if rec.Code != 201 {
    c.Error("Did not update map correctly.")
  } else if result["map"].Id == 0 {
    c.Error("Did not return a map.")
  } else if result["map"].Id != s.mapId {
    c.Error("Did not return the correct map.")
  } else if result["map"].Structure != new_structure {
    c.Error("Did not save new structure value.")
  } else if result["map"].SortBy != new_sort_by {
    c.Error("Did not save new sort_by value.")
  }
}



//
//  Testing helpers
//  -> most are taken from the gocraft/web project
//
func newTestRequest(method, path string, body io.Reader) (*http.Request, *httptest.ResponseRecorder) {
  request, _ := http.NewRequest(method, path, body)
  recorder := httptest.NewRecorder()
  return request, recorder
}


func responseBodyToString(rec *httptest.ResponseRecorder) string {
  return strings.TrimSpace(string(rec.Body.Bytes()))
}


func emptyBuffer() io.Reader {
  return bytes.NewBuffer([]byte{})
}


func setAuthorizationHeader(req *http.Request, s *MySuite) {
  req.Header.Set("Authorization", "Bearer " + s.userAuthToken)
}
