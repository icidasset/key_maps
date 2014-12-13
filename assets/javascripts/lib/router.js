K.Router.map(function() {
  this.route("map", { path: "/:slug" });
});


K.ApplicationRoute = Ember.Route.extend({
  model: function() {
    return this.store.find("map");
  }
});


K.MapRoute = Ember.Route.extend({
  model: function(params) {
    return this.store.filter("map", function(m) {
      return m.get("slug") == params.slug;
    });
  }
});
