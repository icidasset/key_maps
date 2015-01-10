K.InputBooleanComponent = Ember.Component.extend({
  classNames: ["input-boolean"],

  is_true: function() {
    return this.get("value") === true;
  }.property("value"),


  is_false: function() {
    return this.get("value") === false;
  }.property("value"),


  actions: {

    activate: function(bool) {
      this.set("value", bool);
    }

  }

});
