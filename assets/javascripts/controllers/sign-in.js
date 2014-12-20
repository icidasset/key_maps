K.SignInController = Ember.Controller.extend(SimpleAuth.LoginControllerMixin, {
  authenticator: "authenticator:custom",
  errors: [],


  actions: {
    sign_in: function() {
      var data = this.getProperties("identification", "password");
      var controller = this;
      var err;

      this.set("password", null);

      // err
      err = function(error) {
        controller.set("errors", [error]);
      };

      // authentication process
      this.get("session")
        .authenticate(this.authenticator, data)
        .then(null, err);
    }
  }
});
