K.SignUpController = Ember.Controller.extend({
  model: function() {
    return {};
  },

  actions: {

    submit: function() {
      console.log(this.get("model"));

      // Ember.$.ajax({
      //   url: this.serverTokenEndpoint,
      //   type: "POST",
      //   data: {  },
      //   dataType: "json",
      //   beforeSend: function(xhr, settings) {
      //     xhr.setRequestHeader("Accept", settings.accepts.json);
      //   }
      // });
    }

  }

});
