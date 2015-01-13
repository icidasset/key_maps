K.MapItemDataComponent = Ember.Component.extend({
  classNames: ["form__map-item-data", "row-wrapper"],


  on_init: function() {
    console.log("<render map-item>");

    this.set(
      "item.structure_data_clone",
      $.extend({}, this.get("item.structure_data"))
    );
  }.on("init"),


  //
  //  Properties
  //
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

      parent_controller.deleted_map_items.push(item);
      parent_controller.get("model").removeObject(item);

      item.deleteRecord();
    }

  }

});
