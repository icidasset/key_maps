Ember.Application.initializer({
  name: "loading",

  initialize: function(container, application) {
    document.body.classList.remove("is-loading");
  }
});
