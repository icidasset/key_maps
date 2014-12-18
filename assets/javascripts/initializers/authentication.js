Ember.Application.initializer({
  name: "authentication",
  before: "simple-auth",

  initialize: function(container, application) {
    container.register("authenticator:custom", K.Authenticator);
    container.register("authorizer:custom", K.Authorizer);
  }
});
