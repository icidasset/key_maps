K.MapManageController = Ember.Controller.extend({
  needs: ["application", "map", "mapIndex"],
  self_destruct_confirmation: false,


  set_default_self_destruct: function() {
    this.set("self_destruct_confirmation", false);
  }.on("didInsertElement").observes("model"),


  //
  //  Actions
  //
  actions: {

    self_destruct: function() {
      this.set("self_destruct_confirmation", true);
    },

    self_destruct_confirmation: function() {
      var model = this.get("model");
      if (model) model.destroyRecord();

      // reset map selector
      var header_component = this.get(
        "controllers.application.header_component"
      );

      header_component.setProperties({
        map_selector_is_idle: true,
        map_selector_value: ""
      });

      // redirect to index
      this.transitionToRoute("index");
    },

    import_data: function() {
      var text = this.get("import_data_text");
      var controller = this;
      var parsed_json;

      try {
        parsed_json = JSON.parse(text);
      } catch (err) {
        this.wuphf.warning("<i class='help'></i> Invalid JSON");
        return;
      }

      if (Object.prototype.toString.call(parsed_json) != "[object Array]") {
        this.wuphf.warning("<i class='help'></i> Invalid JSON");
        return;
      }

      if (parsed_json) {
        this.wuphf.success("<i class='check'></i> Imported");

        parsed_json.forEach(function(x) {
          controller.get("controllers.mapIndex").add_new(x);
        });
      }
    }

  }

});
