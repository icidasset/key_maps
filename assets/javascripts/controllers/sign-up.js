K.SignUpController = Ember.Controller.extend({
  errors: [],


  ERROR_MESSAGES: {
    "email_presence" : "Given email is not valid",
    "password_presence" : "Given password is not valid" +
                          ", it must be at least 5 characters long",
    "password_not_confirmed" : "Password confirmation is invalid"
  },


  set_errors: function(errors) {
    var EM = this.ERROR_MESSAGES;
    var messages = [];

    errors.forEach(function(err) {
      messages.push(EM[err]);
    });

    this.set("errors", messages);
  },


  actions: {

    submit: function() {
      var _this = this;
      var err = function(x) { _this.set_errors.call(_this, x); };
      var m = this.get("model");
      var e = m.get("email");
      var p = m.get("password");
      var pc = m.get("password_confirmation");

      // to lowercase
      e = e.toLowerCase();
      m.set("email", e);

      // validation
      if (!e || e.length < 5) return err(["email_presence"]);
      if (!p || p.length < 5 || !pc) return err(["password_presence"]);
      if (p !== pc) return err(["password_not_confirmed"]);
      else err([]); // clear errors

      // create model
      return m.save().then(function(response) {
        _this.get("session").authenticate(
          "authenticator:custom",
          { token: response.get("token") }
        );
      }, function(xhr) {
        // xhr.responseJSON.error
        _this.set("errors", ["Email is already taken."]);
      });
    }

  }

});
