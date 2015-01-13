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
    var structure_data_clone = this.get("targetObject.item.structure_data_clone");
    var key = this.get("key");

    // getter
    if (arguments.length === 1) {
      return structure_data_clone[key];

    // setter
    } else {
      structure_data_clone[key] = val;
      return val;

    }

  }.property("key", "type")

});
