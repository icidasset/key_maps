package api

import (
  "encoding/json"
  "github.com/gocraft/web"
)


type BaseContext struct {}


type Context struct {
  *BaseContext

  User User

  UserAuthFormData UserAuthFormData
  MapFormData MapFormData
  MapItemFormData MapItemFormData
}


//
//  Errors
//
type FormattedError struct {
  Error string `json:"error"`
}


func FormatError(err error) FormattedError {
  return FormattedError{ Error: err.Error() }
}


//
//  Rendering
//
func RenderJSON(rw web.ResponseWriter, code int, item interface{}) {
  rw.Header().Set("Content-Type", "application/json")

  // to json
  j, err := json.Marshal(item)

  if err != nil {
    j, _ = json.Marshal(FormatError(err))
    rw.WriteHeader(500)
  } else {
    rw.WriteHeader(code)
  }

  // write to response
  rw.Write(j)
}
