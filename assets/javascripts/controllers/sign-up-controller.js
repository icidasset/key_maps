K.SignUpController = Ember.Controller.extend({
  model: function() {
    return {};
  },

  actions: {

    submit: function() {
      console.log(this.get("model"));
    }

  }

});
