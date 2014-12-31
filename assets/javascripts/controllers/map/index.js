K.MapIndexController = Ember.Controller.extend({
  fullWidthTypes: ["text"],
  destroyedMapItems: [],


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

    return all;
  }.property("model.structure"),


  keys: function() {
    return JSON.parse(this.get("model.structure"));
  }.property("model.structure"),


  actions: {

    add: function() {
      var controller = this;

      this.get("model.map_items").then(function() {
        controller.get("model.map_items").addObject(
          controller.store.createRecord("map_item", {})
        );
      });
    },


    save: function() {
      var destroyed_items = this.destroyedMapItems;

      this.get("model.map_items").forEach(function(mi) {
        if (mi.get("isDirty")) mi.save();
      });

      destroyed_items.forEach(function(d) {
        d.save();
      });

      destroyed_items.length = 0;
    }

  }
});