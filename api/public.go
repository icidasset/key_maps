package api

import (
  "encoding/base64"
  "encoding/json"
  "github.com/go-martini/martini"
  "github.com/icidasset/key-maps/db"
  "github.com/martini-contrib/render"
  "net/http"
  "strconv"
)


func Public__Show(params martini.Params, r render.Render, w http.ResponseWriter) {
  data, err := base64.StdEncoding.DecodeString(params["hash"])
  str := string(data[:])

  // params
  map_id, _ := strconv.ParseInt(str, 10, 0)

  // map items
  map_items := []MapItem{}

  db.Inst().Select(
    &map_items,
    "SELECT structure_data FROM map_items WHERE map_id = $1",
    map_id,
  )

  // collection
  collection := make([]map[string]interface{}, 0)

  for _, m := range map_items {
    c := make(map[string]interface{})
    err = json.Unmarshal([]byte(m.StructureData), &c)
    collection = append(collection, c)
  }

  // render
  if err != nil {
    r.JSON(500, err.Error())
  } else {
    r.JSON(200, collection)
  }
}
