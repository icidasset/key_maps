K.MapIndexController = Ember.ArrayController.extend(DebouncedPropertiesMixin, {
  needs: ["map"],

  full_width_types: ["text"],
  deleted_map_items: [],
  sorted_model_with_numbers: null,

  eachStructureProperty: null,
  eachStructurePropertyDelay: 250,
  debouncedProperties: ["eachStructureProperty"],

  // aliases
  keys: Ember.computed.alias("controllers.map.keys"),
  has_keys: Ember.computed.alias("controllers.map.has_keys"),


  //
  //  Observers
  //
  model_observer: function() {
    var timestamp = (new Date()).getTime();
    this.set("eachStructureProperty", timestamp);
  }.observes("keys", "model.[]"),


  each_structure_property_observer: function() {
    Ember.run.once(this, "set_sorted_model_with_numbers");
  }.observes("debouncedEachStructureProperty"),


  makeNewItemWhenThereIsNone: function() {
    var controller = this;

    if (this.get("model.length") === 0 && this.get("hasKeys")) {
      controller.get("controllers.map.model.map_items").then(function() {
        controller.get("controllers.map.model.map_items").addObject(
          controller.store.createRecord("map_item", {})
        );

        controller.send("resetModel");
      });
    }
  }.observes("setSortedModelWithNumbers"),


  //
  //  Properties
  //
  has_data: function() {
    return this.get("sorted_model_with_numbers") !== null;
  }.property("sorted_model_with_numbers"),


  set_sorted_model_with_numbers: function() {
    var items = this.get("model");
    var keys = this.get("keys");

    var sort_by = (
      this.get("controllers.map.model.sort_by") ||
      (keys[0] ? keys[0].key : null)
    );

    items = items.filter(function(m) {
      return !m.get("isDeleted");
    });

    items = items.sort(function(a, b) {
      var a_struct = a.get("structure_data");
      var b_struct = b.get("structure_data");

      a_struct = a_struct ? JSON.parse(a_struct) : null;
      b_struct = b_struct ? JSON.parse(b_struct) : null;

      a_struct = a_struct && sort_by ? a_struct[sort_by] || "" : "";
      b_struct = b_struct && sort_by ? b_struct[sort_by] || "" : "";

      return a_struct.localeCompare(b_struct);
    });

    items.forEach(function(m, idx) {
      m.set("row_number", idx + 1);
    });

    this.set("sorted_model_with_numbers", items);
  },


  struct: function() {
    var keys = this.get("keys");
    var fwt = this.get("full_width_types");
    var full = [];
    var all = [];

    keys.forEach(function(k) {
      var l = all.length === 0 ? undefined : all[all.length - 1];

      if (fwt.contains(k.type)) {
        full.push(k);
      } else {
        if (l === undefined || l.length >= 2) {
          l = [];
          all.push(l);
        }

        l.push(k);
      }
    });

    if (full.length > 0) {
      all.push(full);
    }

    all.forEach(function(a) {
      a.has_one_item = (a.length === 1);
    });

    return all;
  }.property("keys"),


  //
  //  Other
  //
  clean_up_data: function(item, keys) {
    var data = JSON.parse(item.get("structure_data"));
    var data_keys = Object.keys(data);
    var changed_structure = false;

    for (var i=0, j=data_keys.length; i<j; ++i) {
      var key = data_keys[i];
      if (keys.indexOf(key) === -1) {
        delete data[key];
        changed_structure = true;
      }
    }

    if (changed_structure) {
      item.set("structure_data", JSON.stringify(data));
    }
  },


  //
  //  Actions
  //
  actions: {

    add: function() {
      var controller = this;

      controller.get("controllers.map.model.map_items").then(function() {
        controller.get("controllers.map.model.map_items").addObject(
          controller.store.createRecord("map_item", {})
        );

        controller.send("resetModel");
      });
    },


    save: function() {
      var controller = this;

      Ember.run(function() {
        var promises = [];
        var deleted_items = controller.deleted_map_items;
        var keys = controller.get("keys").map(function(k) {
          return k.key;
        });

        // persist deleted items
        deleted_items.forEach(function(d) {
          promises.push(d.save());
        });

        deleted_items.length = 0;

        // clean up data and save modified items
        controller.get("model").forEach(function(item) {
          controller.clean_up_data(item, keys);
          if (item.get("isDirty")) promises.push(item.save());
        });

        // reset model when all requests are done
        Ember.RSVP.all(promises).then(function() {
          controller.send("resetModel");
        });
      });

      // woof
      this.wuphf.success("<i class='check'></i> Saved");
    }

  }
});
