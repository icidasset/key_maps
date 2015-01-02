K.MapIndexController = Ember.ArrayController.extend({
  needs: ["map"],

  fullWidthTypes: ["text"],
  deletedMapItems: [],

  sortedModel: Ember.computed.sort("model", function(a, b) {
    var a_struct = a.get("structure_data");
    var b_struct = b.get("structure_data");

    a_struct = a_struct ? JSON.parse(a_struct) : null;
    b_struct = b_struct ? JSON.parse(b_struct) : null;

    a_struct = a_struct ? a_struct.author || "" : "";
    b_struct = b_struct ? b_struct.author || "" : "";

    return a_struct.localeCompare(b_struct);
  }),


  make_new_item_on_init: function() {
    if (this.get("model.length") === 0 && this.get("keys")[0].key) {
      this.get("model").addObject(
        this.store.createRecord("map_item", {})
      );
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
      this.get("model").addObject(
        this.store.createRecord("map_item", {})
      );
    },


    save: function() {
      var deleted_items = this.deletedMapItems;

      this.get("model").forEach(function(mi) {
        if (mi.get("isDirty")) mi.save();
      });

      deleted_items.forEach(function(d) {
        d.save();
      });

      deleted_items.length = 0;
    }

  }
});
