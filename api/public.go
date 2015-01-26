package api

import (
  "encoding/base64"
  "encoding/json"
  "github.com/gocraft/web"
  "github.com/icidasset/key-maps/db"
  "strconv"
  "strings"
)


func (c *Context) Public__Show(rw web.ResponseWriter, req *web.Request) {
  data, err := base64.StdEncoding.DecodeString(req.PathParams["hash"])
  s := strings.Split(string(data[:]), "/")

  // params
  map_id, _ := strconv.ParseInt(s[0], 10, 0)
  slug := s[1]

  // map
  m := Map{}

  db.Inst().Get(
    &m,
    "SELECT * FROM maps WHERE id = $1 AND slug = $2",
    map_id,
    slug,
  )

  // return if error
  if m.Id == 0 {
    RenderJSON(rw, 501, map[string]string{ "error": "Provided map id and slug do not match" })
    return
  }

  // map settings
  map_settings := MapSettings{}
  json.Unmarshal([]byte(m.Settings), &map_settings)

  // map items
  map_items := []MapItem{}
  map_items_query := "SELECT structure_data FROM map_items WHERE map_id = $1"

  if map_settings.SortBy != "" {
    map_items_query = map_items_query +
      " ORDER BY structure_data::json->>'" + map_settings.SortBy + "'"
  }

  err = db.Inst().Select(
    &map_items,
    map_items_query,
    map_id,
  )

  // return if error
  if err != nil {
    RenderJSON(rw, 500, FormatError(err));
    return
  }

  // collection
  collection := make([]map[string]interface{}, 0)

  for _, m := range map_items {
    c := make(map[string]interface{})
    err = json.Unmarshal([]byte(m.StructureData), &c)
    collection = append(collection, c)
  }

  // render
  if err != nil {
    RenderJSON(rw, 500, FormatError(err))
  } else {
    RenderJSON(rw, 200, collection)
  }
}
