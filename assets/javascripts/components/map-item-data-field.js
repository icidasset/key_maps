K.MapItemDataFieldComponent = Ember.Component.extend({
  classNames: ["field"],
  classNameBindings: [
    "is_type_text:is-full-width",
    "is_type_text:has-textarea-height",
    "is_type_boolean:has-normal-height",
    "is_other_type:has-normal-height"
  ],


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
    var values = this.get("targetObject.item.structure_data");
    var key = this.get("key");

    // getter
    if (arguments.length === 1) {
      return values[key];

    // setter
    } else {
      values[key] = val;
      return val;

    }

  }.property("key", "type")

});
