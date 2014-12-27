K.MapItemDataFieldComponent = Ember.Component.extend({
  classNames: ["field"],


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

  }.property("key", "type")

});
