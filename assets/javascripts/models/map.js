var attr = DS.attr;


App.Map = DS.Model.extend({
  name: attr("string"),
  slug: attr("string"),
  structure: attr(),
  created_at: attr(),
  updated_at: attr()
});
