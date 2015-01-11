K.MapIndexController = Ember.ArrayController.extend(DebouncedPropertiesMixin, {
  needs: ["map"],

  full_width_types: ["text"],
  deleted_map_items: [],

  // aliases
  keys: Ember.computed.alias("controllers.map.keys"),
  keys_object: Ember.computed.alias("controllers.map.keys_object"),
  has_keys: Ember.computed.alias("controllers.map.has_keys"),


  //
  //  Observers
  //
  model_observer: function() {
    Ember.run.once(Ember.run.bind(this, this.set_sorted_model));
  }.observes("keys", "model.[]"),


  make_new_item_when_there_is_none: function() {
    if (this.get("model.length") === 0 && this.get("hasKeys")) {
      this.add_new();
    }
  }.observes("sorted_model"),


  //
  //  Properties
  //
  has_data: function() {
    return this.get("sorted_model") !== null;
  }.property("sorted_model"),


  sort_by: function() {
    var keys = this.get("keys");

    return (
      this.get("controllers.map.model.sort_by") ||
      (keys[0] ? keys[0].key : null)
    );
  }.property("controllers.map.model.sort_by", "keys"),


  struct: function() {
    var keys = this.get("keys");
    var fwt = this.get("full_width_types");
    var all = [];

    keys.forEach(function(k) {
      var l = all.length === 0 ? undefined : all[all.length - 1];

      if (fwt.contains(k.type)) {
        all.push([k]);
        all.push([]);
      } else {
        if (l === undefined || l.length >= 2) {
          l = [];
          all.push(l);
        }

        l.push(k);
      }
    });

    all.forEach(function(a) {
      a.has_one_item = (a.length === 1);
    });

    return all;
  }.property("keys"),


  set_sorted_model: function() {
    var items = this.get("model").toArray();
    var sort_by = this.get("sort_by");

    items = items.filter(function(m) {
      return !m.get("isDeleted");
    });

    items = items.sort(function(a, b) {
      var a_struct = a.get("structure_data");
      var b_struct = b.get("structure_data");

      a_struct = a_struct ? JSON.parse(a_struct) : null;
      b_struct = b_struct ? JSON.parse(b_struct) : null;

      a_struct = a_struct && sort_by ? a_struct[sort_by] || "" : "";
      b_struct = b_struct && sort_by ? b_struct[sort_by] || "" : "";

      return a_struct.toString().localeCompare(b_struct.toString());
    });

    this.set("sorted_model", items);
  },


  //
  //  Other
  //
  clean_up_data: function(item, keys) {
    var keys_object = this.get("keys_object");
    var was_not_set = !item.get("structure_data");
    var data;

    if (was_not_set) {
      data = {};

      keys.forEach(function(k) {
        data[k] = null;
      });
    } else {
      data = JSON.parse(item.get("structure_data"));
    }

    var data_keys = Object.keys(data);
    var changed_structure = was_not_set;

    for (var i=0, j=data_keys.length; i<j; ++i) {
      var key = data_keys[i];
      if (keys.indexOf(key) === -1) {
        delete data[key];
        changed_structure = true;
      } else if (keys_object[key] == "number") {
        data[key] = parseFloat(data[key]);
        changed_structure = true;
      }
    }

    if (changed_structure) {
      item.set("structure_data", JSON.stringify(data));
    }
  },


  add_new: function(data) {
    var controller = this;
    var keys_array = Object.keys(this.get("keys_object"));

    data = data || {};
    data = { structure_data: JSON.stringify(data) };

    controller.get("controllers.map.model.map_items").then(function() {
      controller.get("controllers.map.model.map_items").addObject(
        controller.store.createRecord("map_item", data)
      );
    });
  },


  //
  //  Actions
  //
  actions: {

    add: function() {
      this.add_new();
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

        // Ember.RSVP.all(promises).then(function() {});
      });

      // woof
      this.wuphf.success("<i class='check'></i> Saved");
    }

  }
});
