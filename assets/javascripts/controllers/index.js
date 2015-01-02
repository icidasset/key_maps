K.IndexController = Ember.Controller.extend({

  on_insert: function() {
    console.log($(".map-selector__input"));
    $(".map-selector__input input").focus();
  }.on("init")

});
