K.MapItemDataFieldComponent = Ember.Component.extend({
  classNames: ["field"],
  classNameBindings: ["is_type_text:full-width"],


  is_type_text: function() {
    return this.get("type") == "text";
  }.property("type"),


  is_type_boolean: function() {
    return this.get("type") == "boolean";
  }.property("type"),


  is_other_type: function() {
    var t = this.get("type");
    return t != "text" && t != "boolean";
  }.property("type"),


  fieldValue: function(k, val, old_val) {
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
    this.notifyPropertyChange("fieldValue");
  }.observes("targetObject.values")

});
