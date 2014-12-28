K.ApplicationController = Ember.Controller.extend({
  searcher: null,


  setup_searcher: function() {
    var s = this.get("searcher");
    var model = this.get("model");

    if (!model) {
      return;
    }

    if (!s) {
      s = new Sifter();
      this.set("searcher", s);
    }

    s.items = model.map(function(m) {
      return {
        name: m.get("name"),
        slug: m.get("slug")
      };
    });
  }.observes("model.@each.slug")

});
