package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/icidasset/key-maps-api/api"
	"github.com/icidasset/key-maps-api/db"
	"github.com/icidasset/key-maps-api/middleware"
	"github.com/labstack/echo"
	. "gopkg.in/check.v1"
)

// setup gocheck
func Test(t *testing.T) { TestingT(t) }

type MySuite struct {
	router *echo.Echo

	// users
	userAuthToken string

	// maps
	mapId       int
	mapSlug     string
	mapSettings api.MapSettings

	// map items
	mapItemId  int
	mapItemId2 int
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
	s.router = echo.New()
	s.router.Use(middleware.Gzip)
	s.router.Use(middleware.Cors)

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
	(s).testApiUsers(c)
	(s).testApiMaps__Part1(c)
	(s).testApiMapItems__Part1(c)
	(s).testApiPublic(c)
	(s).testApiMaps__Part2(c)
	(s).testApiMapItems__Part2(c)
	(s).testApiMaps__Part3(c)
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

  POST '/users'

  @data --- request-body (application/json)
    { "user": { "email": "TEST@gmail.com", "password": "password" } }

  @return --- response-body (application/json)
    { "token": "some-generated-token" }

*/
func (s *MySuite) testApiUsers__Create(c *C) {
	user := api.UserAuth{Email: "TEST@gmail.com", Password: "password"}
	user_form_data := api.UserAuthFormData{User: user}
	j, _ := json.Marshal(user_form_data)

	// make request
	req, rec := newTestRequest("POST", "/users", bytes.NewBuffer(j))
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

  POST '/users/authenticate'

  @data --- request-body (application/json)
    { "user": { "email": "TEST@gmail.com", "password": "password" } }

  @return --- response-body (application/json)
    { "token": "some-generated-token" }

*/
func (s *MySuite) testApiUsers__Authenticate(c *C) {
	user := api.UserAuth{Email: "TEST@gmail.com", Password: "password"}
	user_form_data := api.UserAuthFormData{User: user}
	j, _ := json.Marshal(user_form_data)

	// make request
	req, rec := newTestRequest("POST", "/users/authenticate", bytes.NewBuffer(j))
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

  GET '/users/verify-token?token=TOKEN_FROM_AUTHENTICATE_REQUEST'

  @data --- query-string-param (string)
    'token'

  @return --- response-body (application/json)
    { "is_valid": boolean }

*/
func (s *MySuite) testApiUsers__VerifyToken(c *C) {
	url := "/users/verify-token?token=" + s.userAuthToken

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
}

func (s *MySuite) testApiMaps__Part3(c *C) {
	(s).testApiMaps__Destroy(c)
}

/*

  POST '/maps'

  @data --- request-body (application/json)
    { "map": { "name": "Quotes", "structure": "[]" } }

  @return --- response-body (application/json)
    { "map": { id: 1, "name": "Quotes", "slug": "quotes", "structure": "[]", created_at: "..." } }

*/
func (s *MySuite) testApiMaps__Create(c *C) {
	m := api.Map{Name: "Quotes", Structure: "[]"}
	m_form_data := api.MapFormData{Map: m}
	j, _ := json.Marshal(m_form_data)

	// make request
	req, rec := newTestRequest("POST", "/maps", bytes.NewBuffer(j))
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

  PATCH '/maps/:id'

  @data --- request-body (application/json)
    { "map": { "name": "Quotes", "structure": "..." } }

  @return --- response-body (application/json)
    { "map": { id: 1, "name": "Quotes", "slug": "quotes", "structure": "...", created_at: "..." } }

*/
func (s *MySuite) testApiMaps__Update(c *C) {
	new_structure := `[{ "key": "quote", "type": "text" }, { "key": "author", "type": "string" }]`
	new_settings := api.MapSettings{SortBy: "author"}
	new_settings_bytes, _ := json.Marshal(&new_settings)
	new_settings_string := string(new_settings_bytes)

	// make json
	m := api.Map{Name: "Quotes", Structure: new_structure, Settings: new_settings_string}
	m_form_data := api.MapFormData{Map: m}
	j, _ := json.Marshal(m_form_data)

	// make request
	req, rec := newTestRequest("PATCH", "/maps/"+strconv.Itoa(s.mapId), bytes.NewBuffer(j))
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
	} else if result["map"].Settings != new_settings_string {
		c.Error("Did not save new sort_by value.")
	} else {
		s.mapSettings = new_settings
	}
}

/*

  GET '/maps/:id'

  @return --- response-body (application/json)
    { "map": { id: 1, "name": "Quotes", "slug": "quotes", "structure": "...", created_at: "..." }

*/
func (s *MySuite) testApiMaps__Show(c *C) {
	req, rec := newTestRequest("GET", "/maps/"+strconv.Itoa(s.mapId), emptyBuffer())
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

  GET '/maps'

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
	req, rec := newTestRequest("GET", "/maps", emptyBuffer())
	setAuthorizationHeader(req, s)
	s.router.ServeHTTP(rec, req)

	// parse gzip
	response := rec.Body.Bytes()

	// parse json from response
	result := api.MapIndex{}
	json.Unmarshal(response, &result)

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

  DELETE '/maps/:id'

*/
func (s *MySuite) testApiMaps__Destroy(c *C) {
	req, rec := newTestRequest("DELETE", "/maps/"+strconv.Itoa(s.mapId), emptyBuffer())
	setAuthorizationHeader(req, s)
	s.router.ServeHTTP(rec, req)

	// validate
	if rec.Code != 204 {
		c.Error("Did not delete map successfully.")

	} else {
		req, rec = newTestRequest("GET", "/maps/"+strconv.Itoa(s.mapId), emptyBuffer())
		setAuthorizationHeader(req, s)
		s.router.ServeHTTP(rec, req)

		if rec.Code != 404 {
			c.Error("Map was not deleted.")
		}

		mi := api.MapItem{}
		db.Inst().Get(
			&mi,
			`SELECT * FROM map_items WHERE id = $1`,
			s.mapItemId2,
		)

		if mi.Id != 0 {
			c.Error("Related map items were not deleted.")
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
	(s).testApiMapItems__Create_2(c)
}

/*

  POST '/map_items'

  @data --- request-body (application/json)
    { "map_item": { "structure_data": "INSERT_JSON_HERE", "map_id": "1" } }

  @return --- request-body (application/json)
    { "map_item": { "id": 1, "structure_data": "INSERT_JSON_HERE", "map_id": "1" } }

*/
func (s *MySuite) testApiMapItems__Create(c *C) {
	m := api.MapItem{StructureData: `{ "author": "Author", "quote": "Quote" }`, MapId: s.mapId}
	m_form_data := api.MapItemFormData{MapItem: m}
	j, _ := json.Marshal(m_form_data)

	// make request
	req, rec := newTestRequest("POST", "/map_items", bytes.NewBuffer(j))
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

// create 2 - which will be used in the maps/delete test
func (s *MySuite) testApiMapItems__Create_2(c *C) {
	m := api.MapItem{StructureData: `{ "author": "Author 2", "quote": "Quote 2" }`, MapId: s.mapId}
	m_form_data := api.MapItemFormData{MapItem: m}
	j, _ := json.Marshal(m_form_data)

	// make request
	req, rec := newTestRequest("POST", "/map_items", bytes.NewBuffer(j))
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
		s.mapItemId2 = result["map_item"].Id
	}
}

/*

  PATCH '/map_items/:id'

  @data --- request-body (application/json)
    { "map_item": { "structure_data": "INSERT_NEW_JSON_HERE" } }

  @return --- request-body (application/json)
    { "map_item": { "id": 1, "structure_data": "INSERT_NEW_JSON_HERE", "map_id": 1 } }

*/
func (s *MySuite) testApiMapItems__Update(c *C) {
	new_structure_data := `{ "author": "Epictetus", "quote": "No great thing is created suddenly." }`

	// make json
	m := api.MapItem{StructureData: new_structure_data}
	m_form_data := api.MapItemFormData{MapItem: m}
	j, _ := json.Marshal(m_form_data)

	// make request
	req, rec := newTestRequest("PATCH", "/map_items/"+strconv.Itoa(s.mapItemId), bytes.NewBuffer(j))
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

  GET '/map_items/:id'

  @return --- request-body (application/json)
    { "map_item": { "id": 1, "structure_data": "INSERT_JSON_HERE", "map_id": 1 } }

*/
func (s *MySuite) testApiMapItems__Show(c *C) {
	req, rec := newTestRequest("GET", "/map_items/"+strconv.Itoa(s.mapItemId), emptyBuffer())
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

  DELETE '/map_items/:id'

*/
func (s *MySuite) testApiMapItems__Destroy(c *C) {
	req, rec := newTestRequest("DELETE", "/map_items/"+strconv.Itoa(s.mapItemId), emptyBuffer())
	setAuthorizationHeader(req, s)
	s.router.ServeHTTP(rec, req)

	// validate
	if rec.Code != 204 {
		c.Error("Did not delete map item successfully.")

	} else {
		req, rec = newTestRequest("GET", "/map_items/"+strconv.Itoa(s.mapItemId), emptyBuffer())
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
//  -> Also tests the gzip and cors middleware
//
func (s *MySuite) testApiPublic(c *C) {
	str := strconv.Itoa(s.mapItemId) + "/" + s.mapSlug
	hash := base64.StdEncoding.EncodeToString([]byte(str))

	// -> json
	error_obj := map[string]string{}
	collection_obj := make([]map[string]string, 0)

	// make request
	req, rec := newTestRequest("GET", "/public/"+hash, emptyBuffer())
	req.Header.Set("Accept-Encoding", "gzip")
	s.router.ServeHTTP(rec, req)

	// parse gzip
	response := parseGzipResponse(rec)

	// validate
	if rec.Code != 200 {
		json.Unmarshal(response, &error_obj)
		c.Error("api/public - Error - " + error_obj["error"])

	} else {
		json.Unmarshal(response, &collection_obj)

		if collection_obj[0]["author"] != "Epictetus" {
			c.Error("api/public did not return the correct values.")
		} else if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
			c.Error("api/public response did not have cors headers.")
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
	req.Header.Set("Authorization", "Bearer "+s.userAuthToken)
}

func parseGzipResponse(rec *httptest.ResponseRecorder) []byte {
	reader, _ := gzip.NewReader(rec.Body)
	defer reader.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	return buf.Bytes()
}
