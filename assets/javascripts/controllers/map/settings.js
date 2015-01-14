K.MapSettingsController = Ember.Controller.extend({
  needs: ["map"],

  // aliases
  has_keys: Ember.computed.alias("controllers.map.has_keys"),


  //
  //  Observers
  //
  on_init: function() {
    var s = this.get("controllers.map.model.sort_by");
    var keys = this.get("select_keys");

    if (!s && keys.length) {
      s = keys[0].key;
    }

    this.set(
      "sort_by_select_value",
      { key: s, name: s }
    );
  }.on("init").observes("model"),


  //
  //  Properties
  //
  select_keys: function() {
    var keys = this.get("controllers.map.keys");

    return keys.map(function(k) {
      return {
        key: k,
        name: k
      };
    });
  }.property("controllers.map.keys"),


  //
  //  Actions
  //
  actions: {

    save: function() {
      var m = this.get("model");

      // sort by
      m.set("sort_by", this.get("sort_by_select_value").key);

      // save
      m.save();

      // woof
      this.wuphf.success("<i class='check'></i> Saved");
    }

  }

});
