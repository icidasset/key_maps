K.MapIndexController = Ember.Controller.extend({
  fullWidthTypes: ["text"],
  values: {},


  structure_arrays: function() {
    var structure = JSON.parse(this.get("model.structure"));
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

    all.push(full);

    return all;
  }.property("model.structure"),


  actions: {

    add: function() {
      //
    },


    save: function() {
      //
    }

  }
});
