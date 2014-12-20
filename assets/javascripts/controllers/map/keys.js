K.MapKeysController = Ember.Controller.extend({
  structure: [{}],


  copy_structure: function() {
    this.set(
      "structure",
      JSON.parse(this.get("model.structure"))
    );
  }.observes("model"),


  actions: {

    add: function() {
      var s = this.get("structure");
      var c = s.slice(0, s.length);

      c.push({});

      this.set("structure", c);
    },


    save: function() {
      var m = this.get("model");

      m.set("structure", JSON.stringify(this.get("structure")));
      m.save();
    }

  }
});
