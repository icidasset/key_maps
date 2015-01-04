K.MapIndexController = Ember.ArrayController.extend(DebouncedPropertiesMixin, {
  needs: ["map"],

  fullWidthTypes: ["text"],
  deletedMapItems: [],
  sortedModelWithNumbers: null,

  eachStructureProperty: null,
  eachStructurePropertyDelay: 250,
  debouncedProperties: ["eachStructureProperty"],


  hasData: function() {
    return this.get("sortedModelWithNumbers") !== null;
  }.property("sortedModelWithNumbers"),


  hasKeys: function() {
    var k = this.get("keys");
    return k && k.length && Object.keys(k[0]).length;
  }.property("keys"),


  modelObserver: function() {
    var timestamp = (new Date()).getTime();
    this.set("eachStructureProperty", timestamp);
  }.observes("keys", "model.@each.structure_data"),


  eachStructurePropertyObserver: function() {
    Ember.run.once(this, "setSortedModelWithNumbers");
  }.observes("debouncedEachStructureProperty"),


  setSortedModelWithNumbers: function() {
    var s = this.get("model");

    var sort_by = (
      this.get("controllers.map.model.sort_by") ||
      this.get("keys")[0].key
    );

    s = s.filter(function(m) {
      return !m.get("isDeleted");
    });

    s = s.sort(function(a, b) {
      var a_struct = a.get("structure_data");
      var b_struct = b.get("structure_data");

      a_struct = a_struct ? JSON.parse(a_struct) : null;
      b_struct = b_struct ? JSON.parse(b_struct) : null;

      a_struct = a_struct ? a_struct[sort_by] || "" : "";
      b_struct = b_struct ? b_struct[sort_by] || "" : "";

      return a_struct.localeCompare(b_struct);
    });

    s.forEach(function(m, idx) {
      m.set("row_number", idx + 1);
    });

    this.set("sortedModelWithNumbers", s);
  },


  make_new_item_on_init: function() {
    var controller = this;

    if (this.get("model.length") === 0 && this.get("keys")[0].key) {
      controller.get("controllers.map.model.map_items").then(function() {
        controller.get("controllers.map.model.map_items").addObject(
          controller.store.createRecord("map_item", {})
        );

        controller.send("resetModel");
      });
    }
  }.observes("model"),


  struct: function() {
    var structure = this.get("keys");
    var fwt = this.get("fullWidthTypes");
    var full = [];
    var all = [];

    structure.forEach(function(s) {
      var l = all.length === 0 ? undefined : all[all.length - 1];

      if (fwt.contains(s.type)) {
        full.push(s);
      } else {
        if (l === undefined || l.length >= 2) {
          l = [];
          all.push(l);
        }

        l.push(s);
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


  keys: function() {
    return JSON.parse(this.get("controllers.map.model.structure"));
  }.property("controllers.map.model.structure"),



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
        var deleted_items = controller.deletedMapItems;
        var promises = [];

        deleted_items.forEach(function(d) {
          promises.push(d.save());
        });

        deleted_items.length = 0;

        controller.get("model").forEach(function(mi) {
          if (mi.get("isDirty")) promises.push(mi.save());
        });

        Ember.RSVP.all(promises).then(function() {
          controller.send("resetModel");
        });
      });
    }

  }
});
