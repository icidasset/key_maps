K.MapItemDataComponent = Ember.Component.extend({
  classNames: ["form__map-item"],
  values: {},


  setup_model: function() {
    var s = JSON.parse(this.get("item.structure_data") || "{}");
    var keys = this.get("keys");

    keys.forEach(function(k) {
      k = k.key;
      s[k] = s[k] || null;
    });

    this.set("values", s);
  }.observes("keys").on("init"),


  values_changed: function() {
    if (this._state.toLowerCase() == "indom") {
      this.set(
        "item.structure_data",
        JSON.stringify(this.get("values"))
      );
    }
  }.observes("values")

});
