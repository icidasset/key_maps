const MESSAGES = {
  create: "Press enter to create",
  select: "Press enter to select",
  select_or_create: "Press enter to select" +
    " <span>or</span> " +
    "shift-enter to create"
};


K.ApplicationHeaderComponent = Ember.Component.extend({
  tagName: "header",
  classNames: ["mod-application-header"],

  // properties
  map_selector_value: "",
  map_selector_message: MESSAGES.select_or_create,
  map_selector_show_message: false,

  map_match: false,
  map_match_mask: null,
  map_status: null,


  // handlers
  _register: function() {
    this.set("register-as", this);
  }.on("init"),


  // observations
  map_selector_value_changed: function() {
    var val, match, match_mask, name, is_absolute_match, status;
    var fuz = this.get("targetObject.fuzzy_search");

    val = this.get("map_selector_value");
    if (val) match = fuz.search(val)[0];

    if (match) {
      name = match.item.name;

      is_absolute_match = (
        (name.toLowerCase() == val.toLowerCase()) &&
        (match.score === 0)
      );

      match_mask = (
        ("<span>" + val + "</span>") +
        (name.substr(val.length))
      );
    }

    if (is_absolute_match) {
      status = "select";
    } else if (match) {
      status = "select_or_create";
    } else if (val && val.length > 0) {
      status = "create";
    }

    this.setProperties({
      map_match: match,
      map_match_mask: match_mask,
      map_selector_show_message: (val.length > 0),
      map_status: status
    });
  }.observes("map_selector_value"),


  set_map_selector_message: function() {
    var status = this.get("map_status");
    this.set("map_selector_message", MESSAGES[status]);
  }.observes("map_status"),


  // actions
  actions: {

    input_key_up: function(val, e) {
      var match = this.get("map_match");

      if (e.which == 13) {
        switch (this.get("map_status")) {
          case "create":
            this.create_map(val);
            break;

          case "select":
            this.select_map(match.item.slug);
            break;

          case "select_or_create":
            if (e.shiftKey) this.create_map(val);
            else this.select_map(match.item.slug);
            break;

          default:
            this.get("targetObject").transitionToRoute("index");
        }

      }
    }

  },


  // other
  create_map: function(name) {
    var controller = this.get("targetObject");
    var comp = this, new_map;

    new_map = controller.store.createRecord("map", {
      name: name,
      structure: "[{}]"
    });

    new_map.save().then(function() {
      comp.select_map(new_map.get("slug"));
    }, function() {
      controller.transitionToRoute("index");
    });
  },


  select_map: function(slug) {
    var transition = this.get("targetObject").transitionToRoute("map", slug);
    var match;

    match = this.get("map_match");

    if (match) this.set("map_selector_value", match.item.name);
    this.set("map_selector_show_message", false);

    document.activeElement.blur();
  }

});
