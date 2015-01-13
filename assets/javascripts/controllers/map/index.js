K.MapIndexController = Ember.Controller.extend({
  needs: ["map"],

  full_width_types: ["text"],
  deleted_map_items: [],

  // aliases
  keys: Ember.computed.alias("controllers.map.keys"),
  keys_object: Ember.computed.alias("controllers.map.keys_object"),
  has_keys: Ember.computed.alias("controllers.map.has_keys"),

  // TODO: arrayComputed
  // http://emberjs.com/api/#method_arrayComputed
  //
  // flaggedModel: Ember.arrayComputed("model", {
  //   addedItem: function(array, item, changeMeta, instanceMeta) {
  //     console.log(changeMeta.item.id, this.get("halt_model_changes"));
  //     if (!this.get("halt_model_changes")) {
  //       array.insertAt(changeMeta.index, item);
  //     }
  //     return array;
  //   },
  //
  //   removedItem: function(array, item, changeMeta, instanceMeta) {
  //     console.log("removed yo", this.get("halt_model_changes"));
  //     if (!this.get("halt_model_changes")) {
  //       array.removeAt(changeMeta.index, 1);
  //     }
  //     return array;
  //   }
  // }),

  // filtered
  filteredModel: Ember.computed.filterBy("model", "isDeleted", false),

  // sorted
  sortedSortProperties: [],
  sortedModel: Ember.computed.sort("filteredModel", "sortedSortProperties"),


  //
  //  Observers
  //
  // make_new_item_when_there_is_none: function() {
  //   if (this.get("model.length") === 0 && this.get("hasKeys")) {
  //     this.add_new();
  //   }
  // }.observes("sortedModel"),


  sort_by_observer: function() {
    this.set(
      "sortedSortProperties",
      ["structure_data." + this.get("sort_by") + ":asc"]
    );
  }.observes("sort_by").on("init"),


  //
  //  Properties
  //
  has_data: function() {
    return this.get("model") !== null;
  }.property("model"),


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


  //
  //  Other
  //
  clean_up_data: function(item, keys) {
    var keys_object = this.get("keys_object");
    var data = $.extend({}, item.get("structure_data_clone"));

    var data_keys = Object.keys(data);
    var changed_structure = true; // TODO, check for object equality

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
      item.set("structure_data", data);
    }
  },


  add_new: function(data) {
    var controller = this;
    var keys_array = Object.keys(this.get("keys_object"));

    data = data || {};
    data = { structure_data: data };

    controller.get("model").addObject(
      controller.store.createRecord("map_item", data)
    );
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

        // after
        Ember.RSVP.all(promises).then(function() {
          console.log("<save>");
          // TODO: refresh model?
        });
      });

      // woof
      this.wuphf.success("<i class='check'></i> Saved");
    }

  }
});
