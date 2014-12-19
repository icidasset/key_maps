K.Router.map(function() {
  this.route("sign_in", { path: "/sign-in" });
  this.route("sign_up", { path: "/sign-up" });

  // authenticated routes
  this.route("map", { path: "/:slug" });
});


K.ApplicationRoute = Ember.Route.extend(SimpleAuth.ApplicationRouteMixin, {
  model: function() {
    return this.get_model();
  },

  get_model: function() {
    if (!this.get("session.isAuthenticated")) {
      return { maps: [] };
    } else {
      return this.store.find("map");
    }
  },

  actions: {
    reset_model: function() {
      this.set("model", this.get_model());
    }
  }
});



//
//  Sign in/up
//
K.SignUpRoute = Ember.Route.extend(SimpleAuth.UnauthenticatedRouteMixin, {
  model: function() {
    return this.store.createRecord("user");
  }
});



//
//  Authenticated Routes
//
K.IndexRoute = Ember.Route.extend(SimpleAuth.AuthenticatedRouteMixin);


K.MapRoute = Ember.Route.extend(SimpleAuth.AuthenticatedRouteMixin, {
  model: function(params) {
    return this.store.filter("map", function(m) {
      return m.get("slug") == params.slug;
    });
  }
});
