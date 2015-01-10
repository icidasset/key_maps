package api

import (
  "encoding/base64"
  "encoding/json"
  "github.com/go-martini/martini"
  "github.com/icidasset/key-maps/db"
  "github.com/martini-contrib/render"
  "strconv"
  "strings"
)


func Public__Show(params martini.Params, r render.Render) {
  data, err := base64.StdEncoding.DecodeString(params["hash"])
  s := strings.Split(string(data[:]), "/")

  // params
  map_id, _ := strconv.ParseInt(s[0], 10, 0)

  // map items
  map_items := []MapItem{}

  err = db.Inst().Select(
    &map_items,
    "SELECT structure_data FROM map_items WHERE map_id = $1",
    map_id,
  )

  // return if error
  if err != nil {
    r.JSON(500, FormatError(err));
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
    r.JSON(500, FormatError(err))
  } else {
    r.JSON(200, collection)
  }
}
