App.MainHeaderComponent = Ember.Component.extend({
  tagName: "header",
  classNames: ["mod-header"],
  classNameBindings: ["isMapSelectorOpen:open-map-selector"],

  isMapSelectorOpen: false,

  actions: {
    toggleMapSelector: function() {
      this.toggleProperty("isMapSelectorOpen");
    }
  }

});
