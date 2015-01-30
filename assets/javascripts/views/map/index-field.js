K.MapIndexFieldView = Ember.View.extend({
  classNames: "field",


  fieldValue: function(k, val, old_val) {
    var key = this.get("key");
    var key_array;
    var obj;

    // getter
    if (arguments.length === 1) {
      return this.get("item.structure_changed_data." + key) ||
             this.get("item.structure_data." + key);

    // setter
    } else {
      if (this.get("type") == "number") {
        val = parseFloat(val);
      }

      key_array = key.split(".");

      obj = this.down_the_road(this, key_array, 0);
      obj[key_array[key_array.length - 1]] = val;

      return val;

    }

  }.property(),


  down_the_road: function(view, key_array, step) {
    var base = "item.structure_changed_data";
    var key_chain, obj;

    if (step - 1 < 0) {
      obj = view.get(base);

      if (!obj) {
        obj = {};
        view.set(base, obj);
      }

    } else {
      key_chain = base + "." + key_array.slice(0, step).join(".");
      obj = view.get(key_chain);

      if (!obj) {
        obj = {};
        view.set(key_chain, obj);
      }

    }

    if (step + 1 < key_array.length) {
      return view.down_the_road(view, key_array, step + 1);
    } else {
      return obj;
    }
  }

});
