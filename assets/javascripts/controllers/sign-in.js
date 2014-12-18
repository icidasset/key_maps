K.SignInController = Ember.Controller.extend(SimpleAuth.LoginControllerMixin, {
  authenticator: "authenticator:custom"
});
