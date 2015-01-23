package main

import (
  "bytes"
  "encoding/base64"
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
)

/*

  TODO:

  - gzip

*/

// setup gocheck
func Test(t *testing.T) { TestingT(t) }


type MySuite struct {
  router *web.Router

  // users
  userAuthToken string

  // maps
  mapId int
  mapSlug string
  mapSortBy string

  // map items
  mapItemId int
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
  CreateMapItemRoutes(s.router)
  CreatePublicRoutes(s.router)
}


func (s *MySuite) TearDownSuite(c *C) {
  db.Close()

  drop_db := exec.Command("dropdb", "keymaps_test")
  drop_db.Run()
}


func (s *MySuite) TestAll(c *C) {
  (s).testRootHandler(c)
  (s).testApiUsers(c)
  (s).testApiMaps__Part1(c)
  (s).testApiMapItems__Part1(c)
  (s).testApiPublic(c)
  (s).testApiMaps__Part2(c)
  (s).testApiMapItems__Part2(c)
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
func (s *MySuite) testApiMaps__Part1(c *C) {
  (s).testApiMaps__Create(c)
  (s).testApiMaps__Update(c)
  (s).testApiMaps__Show(c)
}


func (s *MySuite) testApiMaps__Part2(c *C) {
  (s).testApiMaps__Index(c)
  (s).testApiMaps__Destroy(c)
}


/*

  POST '/api/maps'

  @data --- request-body (application/json)
    { "map": { "name": "Quotes", "structure": "[]" } }

  @return --- response-body (application/json)
    { "map": { id: 1, "name": "Quotes", "slug": "quotes", "structure": "[]", created_at: "..." } }

*/
func (s *MySuite) testApiMaps__Create(c *C) {
  m := api.Map{ Name: "Quotes", Structure: "[]" }
  m_form_data := api.MapFormData{ Map: m }
  j, _ := json.Marshal(m_form_data)

  // make request
  req, rec := newTestRequest("POST", "/api/maps", bytes.NewBuffer(j))
  req.Header.Set("Content-Type", "application/json")
  setAuthorizationHeader(req, s)
  s.router.ServeHTTP(rec, req)

  // parse json from response
  result := map[string]api.Map{}
  json.Unmarshal(rec.Body.Bytes(), &result)

  // validate
  if rec.Code != 201 {
    c.Error("Did not create map correctly.")
  } else if result["map"].Id == 0 {
    c.Error("Did not return a map.")
  } else if result["map"].Slug != "quotes" {
    c.Error("Incorrect generation of the slug.")
  } else {
    s.mapId = result["map"].Id
    s.mapSlug = result["map"].Slug
  }
}


/*

  PUT '/api/maps/:id'

  @data --- request-body (application/json)
    { "map": { "name": "Quotes", "structure": "..." } }

  @return --- response-body (application/json)
    { "map": { id: 1, "name": "Quotes", "slug": "quotes", "structure": "...", created_at: "..." } }

*/
func (s *MySuite) testApiMaps__Update(c *C) {
  new_structure := `[{ "key": "quote", "type": "text" }, { "key": "author", "type": "string" }]`
  new_sort_by := `author`

  // make json
  m := api.Map{ Name: "Quotes", Structure: new_structure, SortBy: new_sort_by }
  m_form_data := api.MapFormData{ Map: m }
  j, _ := json.Marshal(m_form_data)

  // make request
  req, rec := newTestRequest("PUT", "/api/maps/" + strconv.Itoa(s.mapId), bytes.NewBuffer(j))
  req.Header.Set("Content-Type", "application/json")
  setAuthorizationHeader(req, s)
  s.router.ServeHTTP(rec, req)

  // parse json from response
  result := map[string]api.Map{}
  json.Unmarshal(rec.Body.Bytes(), &result)

  // validate
  if rec.Code != 200 {
    c.Error("Did not update map correctly.")
  } else if result["map"].Id == 0 {
    c.Error("Did not return a map.")
  } else if result["map"].Id != s.mapId {
    c.Error("Did not return the correct map.")
  } else if result["map"].Structure != new_structure {
    c.Error("Did not save new structure value.")
  } else if result["map"].SortBy != new_sort_by {
    c.Error("Did not save new sort_by value.")
  } else {
    s.mapSortBy = new_sort_by
  }
}


/*

  GET '/api/maps/:id'

  @return --- response-body (application/json)
    { "map": { id: 1, "name": "Quotes", "slug": "quotes", "structure": "...", created_at: "..." }

*/
func (s *MySuite) testApiMaps__Show(c *C) {
  req, rec := newTestRequest("GET", "/api/maps/" + strconv.Itoa(s.mapId), emptyBuffer())
  setAuthorizationHeader(req, s)
  s.router.ServeHTTP(rec, req)

  // parse json from response
  result := map[string]api.Map{}
  json.Unmarshal(rec.Body.Bytes(), &result)

  // validate
  if rec.Code != 200 {
    c.Error("Did not retrieve map correctly.")
  } else if result["map"].Id == 0 {
    c.Error("Did not return a map.")
  } else if result["map"].Id != s.mapId {
    c.Error("Did not return the correct map.")
  }
}


/*

  GET '/api/maps'

  @return --- response-body (application/json)
    {
      maps: [
        ... maps from user ...
      ],

      map_items: [
        ... map items related to retrieved maps ...
      ]
    }

*/
func (s *MySuite) testApiMaps__Index(c *C) {
  req, rec := newTestRequest("GET", "/api/maps/", emptyBuffer())
  setAuthorizationHeader(req, s)
  s.router.ServeHTTP(rec, req)

  // parse json from response
  result := api.MapIndex{}
  json.Unmarshal(rec.Body.Bytes(), &result)

  // validate
  if rec.Code != 200 {
    c.Error("Did not retrieve the maps and their map items correctly.")
  } else if len(result.Maps) != 1 {
    c.Error("Did not return the correct amount of maps.")
  } else if len(result.MapItems) != 1 {
    c.Error("Did not return the correct amount of map items.")
  } else if result.Maps[0].Id != s.mapId {
    c.Error("Did not return the correct map.")
  } else if result.MapItems[0].Id != s.mapItemId {
    c.Error("Did not return the correct map item.")
  }
}


/*

  DELETE '/api/maps/:id'

*/
func (s *MySuite) testApiMaps__Destroy(c *C) {
  req, rec := newTestRequest("DELETE", "/api/maps/" + strconv.Itoa(s.mapId), emptyBuffer())
  setAuthorizationHeader(req, s)
  s.router.ServeHTTP(rec, req)

  // validate
  if rec.Code != 204 {
    c.Error("Did not delete map successfully.")

  } else {
    req, rec = newTestRequest("GET", "/api/maps/" + strconv.Itoa(s.mapId), emptyBuffer())
    setAuthorizationHeader(req, s)
    s.router.ServeHTTP(rec, req)

    if rec.Code != 404 {
      c.Error("Map was not deleted.")
    }

  }
}



//
//  API - Map items
//
func (s *MySuite) testApiMapItems__Part1(c *C) {
  (s).testApiMapItems__Create(c)
  (s).testApiMapItems__Update(c)
  (s).testApiMapItems__Show(c)
}

func (s *MySuite) testApiMapItems__Part2(c *C) {
  (s).testApiMapItems__Destroy(c)
}


/*

  POST '/api/map_items'

  @data --- request-body (application/json)
    { "map_item": { "structure_data": "INSERT_JSON_HERE", "map_id": "1" } }

  @return --- request-body (application/json)
    { "map_item": { "id": 1, "structure_data": "INSERT_JSON_HERE", "map_id": "1" } }

*/
func (s *MySuite) testApiMapItems__Create(c *C) {
  m := api.MapItem{ StructureData: `{ "author": "Author", "quote": "Quote" }`, MapId: s.mapId }
  m_form_data := api.MapItemFormData{ MapItem: m }
  j, _ := json.Marshal(m_form_data)

  // make request
  req, rec := newTestRequest("POST", "/api/map_items", bytes.NewBuffer(j))
  req.Header.Set("Content-Type", "application/json")
  setAuthorizationHeader(req, s)
  s.router.ServeHTTP(rec, req)

  // parse json from response
  result := map[string]api.MapItem{}
  json.Unmarshal(rec.Body.Bytes(), &result)

  // validate
  if rec.Code != 201 {
    c.Error("Did not create map item correctly.")
  } else if result["map_item"].Id == 0 {
    c.Error("Did not return a map item.")
  } else {
    s.mapItemId = result["map_item"].Id
  }
}


/*

  PUT '/api/map_items/:id'

  @data --- request-body (application/json)
    { "map_item": { "structure_data": "INSERT_NEW_JSON_HERE" } }

  @return --- request-body (application/json)
    { "map_item": { "id": 1, "structure_data": "INSERT_NEW_JSON_HERE", "map_id": 1 } }

*/
func (s *MySuite) testApiMapItems__Update(c *C) {
  new_structure_data := `{ "author": "Epictetus", "quote": "No great thing is created suddenly." }`

  // make json
  m := api.MapItem{ StructureData: new_structure_data }
  m_form_data := api.MapItemFormData{ MapItem: m }
  j, _ := json.Marshal(m_form_data)

  // make request
  req, rec := newTestRequest("PUT", "/api/map_items/" + strconv.Itoa(s.mapItemId), bytes.NewBuffer(j))
  req.Header.Set("Content-Type", "application/json")
  setAuthorizationHeader(req, s)
  s.router.ServeHTTP(rec, req)

  // parse json from response
  result := map[string]api.MapItem{}
  json.Unmarshal(rec.Body.Bytes(), &result)

  // validate
  if rec.Code != 200 {
    c.Error("Did not update map item correctly.")
  } else if result["map_item"].Id == 0 {
    c.Error("Did not return a map item.")
  } else if result["map_item"].Id != s.mapItemId {
    c.Error("Did not return the correct map item.")
  } else if result["map_item"].StructureData != new_structure_data {
    c.Error("Did not save new structure-data value.")
  }
}


/*

  GET '/api/map_items/:id'

  @return --- request-body (application/json)
    { "map_item": { "id": 1, "structure_data": "INSERT_JSON_HERE", "map_id": 1 } }

*/
func (s *MySuite) testApiMapItems__Show(c *C) {
  req, rec := newTestRequest("GET", "/api/map_items/" + strconv.Itoa(s.mapItemId), emptyBuffer())
  setAuthorizationHeader(req, s)
  s.router.ServeHTTP(rec, req)

  // parse json from response
  result := map[string]api.MapItem{}
  json.Unmarshal(rec.Body.Bytes(), &result)

  // validate
  if rec.Code != 200 {
    c.Error("Did not retrieve map item correctly.")
  } else if result["map_item"].Id == 0 {
    c.Error("Did not return a map item.")
  } else if result["map_item"].Id != s.mapItemId {
    c.Error("Did not return the correct map item.")
  }
}


/*

  DELETE '/api/map_items/:id'

*/
func (s *MySuite) testApiMapItems__Destroy(c *C) {
  req, rec := newTestRequest("DELETE", "/api/map_items/" + strconv.Itoa(s.mapItemId), emptyBuffer())
  setAuthorizationHeader(req, s)
  s.router.ServeHTTP(rec, req)

  // validate
  if rec.Code != 204 {
    c.Error("Did not delete map item successfully.")

  } else {
    req, rec = newTestRequest("GET", "/api/map_items/" + strconv.Itoa(s.mapItemId), emptyBuffer())
    setAuthorizationHeader(req, s)
    s.router.ServeHTTP(rec, req)

    if rec.Code != 404 {
      c.Error("Map item was not deleted.")
    }

  }
}



//
//  API - Public
//
func (s *MySuite) testApiPublic(c *C) {
  str := strconv.Itoa(s.mapItemId) + "/" + s.mapSlug
  hash := base64.StdEncoding.EncodeToString([]byte(str))

  // -> json
  error_obj := map[string]string{}
  collection_obj := make([]map[string]string, 0)

  // make request
  req, rec := newTestRequest("GET", "/api/public/" + hash, emptyBuffer())
  s.router.ServeHTTP(rec, req)

  // validate
  if rec.Code != 200 {
    json.Unmarshal(rec.Body.Bytes(), &error_obj)
    c.Error("api/public - Error - " + error_obj["error"])

  } else {
    json.Unmarshal(rec.Body.Bytes(), &collection_obj)

    if collection_obj[0]["author"] != "Epictetus" {
      c.Error("api/public did not return the correct values.")
    }

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
