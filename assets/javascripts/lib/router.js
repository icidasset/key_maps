K.Router.map(function() {
  this.route("sign_in", { path: "/sign-in" });
  this.route("sign_up", { path: "/sign-up" });
  this.route("sign_out", { path: "/sign-out" });

  // authenticated routes
  this.route("map", { path: "/:slug" });
});


K.ApplicationRoute = Ember.Route.extend(SimpleAuth.ApplicationRouteMixin, {
  model: function() {
    return this.getModel();
  },

  getModel: function() {
    if (this.get("session.isAuthenticated")) {
      return this.store.find("map");
    }
  },

  actions: {
    sessionAuthenticationSucceeded: function() {
      var route = this;

      this.getModel().then(function(model) {
        route.controller.set("model", model);
      });

      this._super();
    },

    sessionInvalidationSucceeded: function() {
      this.controller.set("model", null);
      this._super();
    }
  }
});



//
//  Sign in/up/out
//
K.SignInRoute = Ember.Route.extend({
  setupController: function(controller, model) {
    controller.set("errors", null);
  }
});


K.SignUpRoute = Ember.Route.extend(SimpleAuth.UnauthenticatedRouteMixin, {
  model: function() {
    return this.store.createRecord("user");
  }
});


K.SignOutRoute = Ember.Route.extend(SimpleAuth.AuthenticatedRouteMixin, {
  beforeModel: function(transition) {
    transition.abort();
    this.send("invalidateSession");
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
