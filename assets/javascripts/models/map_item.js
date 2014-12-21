var attr = DS.attr;


K.MapItem = DS.Model.extend({
  structure_data: attr(),
  created_at: attr(),
  updated_at: attr(),

  map: DS.belongsTo("map")
});
