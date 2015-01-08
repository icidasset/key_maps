K.XWuphfComponent = EmberWuphf.XWuphfComponent;
K.XWuphfMessageComponent = EmberWuphf.XWuphfMessageComponent;


Ember.Application.initializer({
  name: "wuphf",

  initialize: function(container, application) {
    application.register("wuphf:main", EmberWuphf.Service);

    application.inject("controller", "wuphf", "wuphf:main");
    application.inject("component", "wuphf", "wuphf:main");
    application.inject("route", "wuphf", "wuphf:main");
  }
});
