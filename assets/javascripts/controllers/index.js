K.IndexController = Ember.Controller.extend({
  needs: ["application"],

  map_count: function() {
    var n;

    if (this.get("session.isAuthenticated")) {
      n = this.get("store").all("map").get("content").length;
    } else {
      n = 0;
    }

    return n.toString() + " map" + (n !== 1 ? "s" : "");
  }.property("session.isAuthenticated", "controllers.application.model.[]"),


  item_count: function() {
    var n;

    if (this.get("session.isAuthenticated")) {
      n = this.get("store").all("map_item").get("content").length;
    } else {
      n = 0;
    }

    return n.toString() + " map item" + (n !== 1 ? "s" : "");
  }.property("session.isAuthenticated", "controllers.application.model.[]")

});
