K.MapItemDataComponent = Ember.Component.extend({
  classNames: ["form__map-item-data"],
  values: {},


  on_did_insert_element: function() {
    this.addObserver("item.structure_data", this, "setup_model");
    this.notifyPropertyChange("item.structure_data");
  }.on("didInsertElement"),


  setup_model: function() {
    var s = JSON.parse(this.get("item.structure_data") || "{}");
    var keys = this.get("keys");

    if (Object.keys(s).length === 0 &&
        Object.keys(this.get("values")).length === 0) {
      return;
    } else {
      this.removeObserver("item.structure_data", this, "setup_model");
    }

    keys.forEach(function(k) {
      k = k.key;
      s[k] = s[k] || null;
    });

    this.set("values", s);
  },


  values_changed: function() {
    if (this._state.toLowerCase() == "indom") {
      this.set(
        "item.structure_data",
        JSON.stringify(this.get("values"))
      );
    }
  }.observes("values"),


  number: function() {
    return this.get("item.row_number");
  }.property("item.row_number"),



  //
  //  Actions
  //
  actions: {

    destroy: function() {
      var parent_controller = this.get("targetObject");
      var item = this.get("item");

      parent_controller.deleted_map_items.push(item);
      parent_controller.get("model").removeObject(item);

      item.deleteRecord();
    }

  }

});
