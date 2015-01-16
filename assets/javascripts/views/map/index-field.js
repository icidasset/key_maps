K.MapIndexFieldView = Ember.View.extend({
  classNames: "field",


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
