K.MapKeysController = Ember.Controller.extend({
  needs: ["map"],

  structure: [],
  reformatted_structure: [],


  //
  //  Observers
  //
  copy_structure: function() {
    this.set(
      "structure",
      JSON.parse(this.get("model.structure"))
    );
  }.observes("model"),


  reformat_structure: function() {
    var structure = this.get("structure") || [];
    var types = this.get("controllers.map.types");

    var reformatted = structure.map(function(s) {
      var type;

      types.forEach(function(t) {
        if (s.type == t.value) {
          type = t;
        }
      });

      return { key: s.key, type: type };
    });

    if (reformatted.length === 0) {
      reformatted.push({});
    }

    this.set("reformatted_structure", reformatted);
  }.observes("structure"),


  //
  //  Other
  //
  clean_structure: function(structure) {
    var c = [];

    structure.forEach(function(s) {
      if (s.key && s.key.length > 0 && s.type && s.type.value) {
        c.push({ key: s.key, type: s.type.value });
      }
    });

    c = c.sortBy("key");

    return c;
  },



  //
  //  Actions
  //
  actions: {

    add: function() {
      var s = this.get("reformatted_structure");
      var c = s.slice(0, s.length);

      c.push({});

      this.set("reformatted_structure", c);
    },


    save: function() {
      var m = this.get("model");
      var s = this.clean_structure(this.get("reformatted_structure"));

      this.set("structure", s);
      m.set("structure", JSON.stringify(s));
      m.save().then(function() {
        $(document.activeElement).filter("button").blur();
      });
    }

  }

});
