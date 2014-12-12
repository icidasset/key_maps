App.Router.map(function() {
  this.route("map", { path: "/:map_slug" });
});


App.ApplicationRoute = Ember.Route.extend({
  model: function() {
    this.store.find("map");
  }
});


App.MapRoute = Ember.Route.extend({});
