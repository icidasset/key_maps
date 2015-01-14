var attr = DS.attr;


K.Map = DS.Model.extend({
  name: attr("string"),
  slug: attr("string"),
  structure: attr("json_array"),
  sort_by: attr(),
  created_at: attr(),
  updated_at: attr(),

  map_items: DS.hasMany("map_item", { async: true })
});
