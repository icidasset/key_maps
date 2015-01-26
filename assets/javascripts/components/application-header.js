const MESSAGES = {
  create: `Press enter to create`,
  select: `Press enter to select`,
  select_or_create: `Press enter to select
    <span>or</span>
    shift-enter to create`,
  _default: `Type to select or create a map`
};


K.ApplicationHeaderComponent = Ember.Component.extend(Ember.Validations.Mixin, {
  tagName: "header",
  classNames: ["mod-application-header"],

  // properties
  map_selector_value: "",
  map_selector_message: MESSAGES._default,
  map_selector_is_idle: true,

  map_match: false,
  map_match_mask: null,
  map_status: null,


  // status
  has_alert_shown: false,


  // validations
  validations: {
    map_selector_value: {
      length: {
        minimum: 3,
        tokenizer: function(val) {
          return val.replace(/[\W_]+/g, "");
        }
      }
    }
  },


  //
  //  Handlers
  //
  _register: function() {
    this.set("register-as", this);
  }.on("init"),


  _init_element: function() {
    var t = new BareTooltip(
      this.get("element"), {
      trigger_type: "click",
      delegate_selector: ".application-header__menu-trigger"
    });

    t.move_tooltip = function(e) {
      var $tooltip = this.state.$tooltip_element,
          t = e.currentTarget,
          tb = t.getBoundingClientRect(),
          tw = $tooltip.outerWidth();

      var left = tb.left + (tb.width / 2) - (tw / 2) + 2,
          top = tb.top + tb.height + 16;

      if (left < 24) {
        left = tb.left + (tb.width / 2) - 24 + 1;
        $tooltip.addClass("is-left");
      } else if (left + tw > window.innerWidth - 24) {
        left = tb.left + (tb.width / 2) - tw + 24;
        $tooltip.addClass("is-right");
      }

      $tooltip.css({
        left: left,
        top: top
      });
    };

    t.setup();

    this.tooltip = t;
  }.on("didInsertElement"),


  _destroy_element: function() {
    this.tooltip.self_destruct();
    this.tooltip = null;
  }.on("willDestroyElement"),


  //
  //  Observers
  //
  map_selector_value_changed: function() {
    var val, match, result, partial, match_mask, name, is_absolute_match, status;
    var searcher = this.get("targetObject.searcher");

    val = this.get("map_selector_value");

    if (val) {
      result = searcher.search(val, {
        fields: ["name"],
        sort: [{ field: "name", direction: "asc" }],
        limit: 1
      });
    }

    if (result && result.items[0]) {
      match = searcher.items[result.items[0].id];
      partial = match.name.substr(0, result.query.length);

      if (result.query.toLowerCase() != partial.toLowerCase()) {
        match = null;
      }
    }

    if (match) {
      name = match.name;

      is_absolute_match = (
        (name.toLowerCase() == val.toLowerCase()) &&
        (match.score === 0)
      );

      match_mask = (
        ("<span>" + val + "</span>") +
        (name.substr(val.length))
      ).replace(/ /g, "&#160;");
    }

    if (is_absolute_match) {
      status = "select";
    } else if (match) {
      status = "select_or_create";
    } else if (val && val.length > 0) {
      status = "create";
    } else {
      status = "_default";
    }

    this.setProperties({
      map_match: match,
      map_match_mask: match_mask,
      map_status: status
    });
  }.observes("map_selector_value"),


  set_map_selector_message: function() {
    var status = this.get("map_status");
    var is_idle = this.get("map_selector_is_idle");

    this.set(
      "map_selector_message",
      MESSAGES[is_idle ? "_default" : status]
    );
  }.observes(
    "map_status",
    "map_selector_is_idle"
  ),


  //
  //  Other
  //
  create_map: function(name) {
    var controller = this.get("targetObject");
    var comp = this, new_map;

    if (!comp.get("isValid")) {
      comp.set("has_alert_shown", true);

      alert(
        "You must give a valid map name. " +
        "It must contain at least 3 alphanumeric characters."
      );

      setTimeout(function() {
        comp.set("has_alert_shown", false);
      }, 250);

      return;
    }

    comp.setProperties({
      map_selector_is_idle: true,
      map_status: "select"
    });

    new_map = controller.store.createRecord("map", {
      name: name,
      structure: []
    });

    new_map.save().then(function() {
      comp.select_map(new_map.get("slug"), true);
    }, function() {
      controller.transitionToRoute("index");
    });
  },


  select_map: function(slug, via_create=false) {
    var route = via_create ? "map.keys" : "map.index";
    var transition = this.get("targetObject").transitionToRoute(route, slug);
    var match = this.get("map_match");

    if (match) this.set("map_selector_value", match.name);
    this.set("map_selector_is_idle", true);

    document.activeElement.blur();
  },


  //
  //  Actions
  //
  actions: {

    input_key_up: function(val, e) {
      var match = this.get("map_match");

      if (e.which == 13 && !this.get("has_alert_shown")) {
        switch (this.get("map_status")) {
          case "create":
            this.create_map(val);
            break;

          case "select":
            this.select_map(match.slug);
            break;

          case "select_or_create":
            if (e.shiftKey) this.create_map(val);
            else this.select_map(match.slug);
            break;

          default:
            this.get("targetObject").transitionToRoute("index");
        }

      } else {
        this.set("map_selector_is_idle", false);

      }
    },


    reset: function() {
      this.setProperties({
        map_selector_is_idle: true,
        map_selector_value: ""
      });
    }

  }

});
