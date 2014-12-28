K.MapKeysController = Ember.Controller.extend({
  structure: [{}],
  types: [
    { val: "string", name: "String" },
    { val: "text", name: "Text" },
    { val: "number", name: "Number" },
    { val: "boolean", name: "Boolean" }
  ],


  copy_structure: function() {
    this.set(
      "structure",
      JSON.parse(this.get("model.structure"))
    );
  }.observes("model"),


  clean_structure: function(structure) {
    var c = [];

    structure.forEach(function(s) {
      if (s.key && s.key.length > 0) {
        c.push(s);
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
      var s = this.clean_structure(this.get("structure"));

      this.set("structure", s);
      m.set("structure", JSON.stringify(s));
      m.save().then(function() {
        $(document.activeElement).filter("button").blur();
      });
    }

  }
});
