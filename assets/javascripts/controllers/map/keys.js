K.MapKeysController = Ember.Controller.extend({
  structure: [{}],
  reformatted_structure: [{}],

  types: [
    { value: "string", name: "String" },
    { value: "text", name: "Text" },
    { value: "number", name: "Number" },
    { value: "boolean", name: "Boolean" }
  ],


  copy_structure: function() {
    this.set(
      "structure",
      JSON.parse(this.get("model.structure"))
    );
  }.observes("model"),


  reformat_structure: function() {
    var structure = this.get("structure");
    var types = this.types;

    var reformatted = structure.map(function(s) {
      var type;

      types.forEach(function(t) {
        if (s.type == t.value) {
          type = t;
        }
      });

      return { key: s.key, type: type };
    });

    this.set("reformatted_structure", reformatted);
  }.observes("structure"),


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


  actions: {

    add: function() {
      var s = this.get("structure");
      var c = s.slice(0, s.length);

      c.push({});

      this.set("structure", c);
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
