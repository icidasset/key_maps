K.ApplicationController = Ember.Controller.extend({
  fuzzy_search: null,


  setup_fuzzy_search: function() {
    var f = this.get("fuzzy_search");

    if (!f) {
      f = new Fuse(null, {
        keys: ["name"],
        includeScore: true,
        maxPatternLength: 64,
        distance: 0,
        threshold: 0
      });

      this.set("fuzzy_search", f);
    }

    f.list = this.get("model").map(function(m) {
      return {
        name: m.get("name"),
        slug: m.get("slug")
      };
    });
  }.observes("model.@each", "model.@each.slug")

});
