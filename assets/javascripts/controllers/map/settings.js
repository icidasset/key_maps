K.MapSettingsController = Ember.Controller.extend({
  needs: ["map"],


  keys: function() {
    var keys = JSON.parse(this.get("controllers.map.model.structure"));
    return keys.map(function(k) {
      return {
        key: k.key,
        name: k.key
      };
    });
  }.property("controllers.map.model.structure"),


  sort_by_observer: function() {
    var k = this.get("controllers.map.model.sort_by");

    this.set(
      "sort_by_select_value",
      { key: k, name: k }
    );
  }.on("init"),


  //
  //  Actions
  //
  actions: {

    save: function() {
      var m = this.get("model");

      m.set("sort_by", this.get("sort_by_select_value").key);

      m.save().then(function() {
        $(document.activeElement).filter("button").blur();
      });
    }

  }

});
