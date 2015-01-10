K.MapController = Ember.Controller.extend({
  needs: "application",

  types: [
    { value: "string", name: "String" },
    { value: "text", name: "Text" },
    { value: "number", name: "Number" },
    { value: "boolean", name: "Boolean" }
  ],


  //
  //  Observers
  //
  pass_map_name_to_header: function() {
    var m;
    var header_component = this.get(
      "controllers.application.header_component"
    );

    // check
    if (!header_component) return;

    // continue
    m = this.get("model");

    if (m) {
      header_component.set("map_selector_value", m.get("name"));
      header_component.set("map_selector_show_message", false);
      document.activeElement.blur();
    }
  }.observes(
    "model",
    "controllers.application.header_component"
  ),


  //
  //  Properties
  //
  has_keys: function() {
    var k = this.get("keys");
    return k && k.length && Object.keys(k[0]).length;
  }.property("keys"),


  keys: function() {
    return JSON.parse(this.get("model.structure"));
  }.property("model.structure"),


  keys_object: function() {
    var o = {};

    this.get("keys").forEach(function(k) {
      o[k.key] = k.type;
    });

    return o;
  }.property("keys"),


  public_url: function() {
    var id = this.get("model.id");
    var slug = this.get("model.slug");
    var base = Base64.encode(id.toString() + "/" + slug);

    return "" +
      window.location.protocol + "//" +
      window.location.host + "/api/public/" + base;
  }.property(
    "model.id",
    "model.slug",
    "model.map_items.[]",
    "model.map_items.@each.structure_data"
  )

});
