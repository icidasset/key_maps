K.MapItemDataFieldComponent = Ember.Component.extend({
  classNames: ["field"],
  classNameBindings: ["is_type_text:full-width"],


  is_type_text: function() {
    return this.get("type") == "text";
  }.property("type"),


  field_value: function(k, val, old_val) {
    var values = this.get("targetObject.values");
    var key = this.get("key");

    // getter
    if (arguments.length === 1) {
      return values[key];

    // setter
    } else {
      values[key] = val;
      this.get("targetObject").notifyPropertyChange("values");
      return val;

    }

  }.property("key", "type"),


  observe_item: function() {
    this.notifyPropertyChange("field_value");
  }.observes("targetObject.item.structure_data")

});
