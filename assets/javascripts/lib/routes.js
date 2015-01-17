K.Router.map(function() {
  this.route("sign_in", { path: "/sign-in" });
  this.route("sign_up", { path: "/sign-up" });
  this.route("sign_out", { path: "/sign-out" });

  // authenticated routes
  this.route("map", { path: "/:slug" }, function() {
    this.route("keys", { path: "/keys" });
    this.route("settings", { path: "/settings" });
    this.route("manage", { path: "/manage" });
  });
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
      if (this.controller) {
        this.controller.set("model", null);
      }

      this._super();
    }
  }
});


K.IndexRoute = Ember.Route.extend({
  afterModel: function() {
    if (this.get("session.isAuthenticated")) {
      var ac = this.controllerFor("application");
      var header_component = ac.get("header_component");
      if (header_component) header_component.send("reset");
    }
  }
});



//
//  Sign in/up/out
//
K.SignInRoute = Ember.Route.extend({
  setupController: function(controller, model) {
    controller.set("errors", null);
  },

  beforeModel: function(transition) {
    if (this.get("session.isAuthenticated")) {
      this.transitionTo("index");
    }
  }
});


K.SignUpRoute = Ember.Route.extend(SimpleAuth.UnauthenticatedRouteMixin, {
  model: function() {
    return this.store.createRecord("user");
  },

  beforeModel: function(transition) {
    if (this.get("session.isAuthenticated")) {
      this.transitionTo("index");
    }
  }
});


K.SignOutRoute = Ember.Route.extend(SimpleAuth.AuthenticatedRouteMixin, {
  beforeModel: function(transition) {
    if (this.get("session.isAuthenticated")) {
      transition.send("invalidateSession");
    } else {
      this.transitionTo("index");
    }
  }
});



//
//  Authenticated Routes
//
K.MapRoute = Ember.Route.extend(SimpleAuth.AuthenticatedRouteMixin, {
  model: function(params) {
    var m = this.getModel(params);
    if (m) return m;
    else return null;
  },

  getModel: function(params) {
    return this.store.all("map").filter(function(m) {
      return m.get("slug") == params.slug;
    })[0];
  }
});


K.MapIndexRoute = Ember.Route.extend(SimpleAuth.AuthenticatedRouteMixin, {
  model: function() {
    return this.getModel();
  },

  getModel: function() {
    var m = this.modelFor("map");
    if (m) return m.get("map_items");
    else return null;
  },

  resetModel: function() {
    this.set("controller.model", this.getModel());
  }
});
