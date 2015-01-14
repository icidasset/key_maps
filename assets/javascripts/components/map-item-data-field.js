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
  }.property("type").readOnly(),


  is_type_boolean: function() {
    return this.get("type") == "boolean";
  }.property("type").readOnly(),


  is_other_type: function() {
    var t = this.get("type");
    return t != "text" && t != "boolean";
  }.property("type").readOnly(),


  fieldValue: function(k, val, old_val) {
    var structure_data = this.get("targetObject.item.structure_data");
    var structure_changed_data = this.get("targetObject.item.structure_changed_data");
    var key = this.get("key");

    // check
    if (!structure_changed_data) {
      this.set("targetObject.item.structure_changed_data", {});
      structure_changed_data = this.get("targetObject.item.structure_changed_data");
    }

    // getter
    if (arguments.length === 1) {
      return structure_changed_data[key] || structure_data[key];

    // setter
    } else {
      structure_changed_data[key] = val;
      return val;

    }

  }.property("key")

});
