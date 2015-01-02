K.MapItemDataComponent = Ember.Component.extend({
  classNames: ["form__map-item-data"],
  values: {},


  setup_model: function() {
    var s = JSON.parse(this.get("item.structure_data") || "{}");
    var keys = this.get("keys");

    keys.forEach(function(k) {
      k = k.key;
      s[k] = s[k] || null;
    });

    this.set("values", s);
  }.observes("item.structure_data").on("init"),


  values_changed: function() {
    if (this._state.toLowerCase() == "indom") {
      this.set(
        "item.structure_data",
        JSON.stringify(this.get("values"))
      );
    }
  }.observes("values"),


  number: function() {
    return this.get("idx") + 1;
  }.property("idx"),



  //
  //  Actions
  //
  actions: {

    destroy: function() {
      var parent_controller = this.get("targetObject");
      var item = this.get("item");

      parent_controller.deletedMapItems.push(item);
      item.deleteRecord();
    }

  }

});
