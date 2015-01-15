K.MapItemDataFieldComponent = Ember.Component.extend({
  classNames: "field",

  did_insert_element: function() {
    var el = this.get("element");

    if (this.get("is_type_text")) {
      el.classList.add("is-full-width");
      el.classList.add("has-textarea-height");
    } else {
      el.classList.add("has-normal-height");
    }
  }.on("didInsertElement"),


  is_type_text: function() {
    return this.get("type") == "text";
  }.property().readOnly(),


  is_type_boolean: function() {
    return this.get("type") == "boolean";
  }.property().readOnly(),


  is_other_type: function() {
    var t = this.get("type");
    return t != "text" && t != "boolean";
  }.property().readOnly(),


  fieldValue: function(k, val, old_val) {
    var structure_data = this.get("item.structure_data");
    var structure_changed_data = this.get("item.structure_changed_data");
    var key = this.get("key");

    // check
    if (!structure_changed_data) {
      this.set("item.structure_changed_data", {});
      structure_changed_data = this.get("item.structure_changed_data");
    }

    // getter
    if (arguments.length === 1) {
      return structure_changed_data[key] || structure_data[key];

    // setter
    } else {
      structure_changed_data[key] = val;
      return val;

    }

  }.property()

});
