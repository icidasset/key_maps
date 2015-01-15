K.MapIndexController = Ember.Controller.extend({
  needs: ["map"],

  full_width_types: ["text"],
  deleted_map_items: [],
  halt_model_changes: false,

  // aliases
  keys: Ember.computed.readOnly("controllers.map.keys"),
  keys_object: Ember.computed.readOnly("controllers.map.keys_object"),
  has_keys: Ember.computed.readOnly("controllers.map.has_keys"),

  // check for halt-model-changes flag
  flaggedModel: Ember.arrayComputed("model", {
    addedItem: function(array, item, changeMeta, instanceMeta) {
      if (!this.get("halt_model_changes")) {
        array.insertAt(changeMeta.index, item);
      }
      return array;
    },

    removedItem: function(array, item, changeMeta, instanceMeta) {
      if (!this.get("halt_model_changes")) {
        array.removeAt(changeMeta.index, 1);
      }
      return array;
    }
  }),

  // filtered
  filteredModel: Ember.computed.filterBy("flaggedModel", "isDeleted", false),

  // sorted
  sortedSortProperties: [],
  sortedModel: Ember.computed.sort("filteredModel", "sortedSortProperties"),


  //
  //  Observers
  //
  make_new_item_when_there_is_none: function() {
    if (this.get("model.length") === 0 && this.get("has_keys")) {
      this.add_new();
    }
  }.observes("model"),


  sort_by_observer: function() {
    this.set(
      "sortedSortProperties",
      ["structure_data." + this.get("sort_by") + ":asc"]
    );
  }.observes("sort_by").on("init"),


  //
  //  Properties
  //
  sort_by: function() {
    return this.get("controllers.map.model.sort_by") || this.get("keys")[0];
  }.property(
    "controllers.map.model.sort_by",
    "keys"
  ).readOnly(),


  struct: function() {
    var structure = this.get("controllers.map.model.structure");
    var fwt = this.get("full_width_types");
    var all = [];

    structure.forEach(function(k) {
      var l = all.length === 0 ? undefined : all[all.length - 1];

      if (fwt.contains(k.type)) {
        all.push([k]);
        all.push([]);
      } else {
        if (l === undefined || l.length >= 2) {
          l = [];
          all.push(l);
        }

        l.push(k);
      }
    });

    all.forEach(function(a) {
      a.has_one_item = (a.length === 1);
    });

    return all;
  }.property(
    "controllers.map.model.structure"
  ).readOnly(),


  item_template: function() {
    var t = [
      '<div class="row-prefix" {{action "destroy"}}>',
        '<span class="row-prefix__title row-prefix__center">',
          '{{#if item.isNew}}NEW{{else}}{{number}}{{/if}}',
        '</span>',
        '<span class="row-prefix__destroy row-prefix__center">',
          '<i class="cross"></i>',
        '</span>',
      '</div>'
    ].join("");

    this.get("struct").forEach(function(s) {
      var row_class = "row " + (s.length === 1 ? "row__with-one-item" : "");

      // <row>
      t = t + '<div class="' + row_class + '">';

      // fields
      s.forEach(function(field) {
        t = t + '{{map-item-data-field key="' + field.key + '" type="' + field.type + '"}}';
      });

      // </row>
      t = t + '</div>';
    });

    return Ember.Handlebars.compile(t);
  }.property("struct").readOnly(),


  //
  //  Other
  //
  clean_up_data: function(item) {
    var keys = this.get("keys");
    var keys_object = this.get("keys_object");
    var structure_data = item.get("structure_data");
    var structure_changed_data = item.get("structure_changed_data");

    var changed_structure = (
      structure_changed_data &&
      Object.keys(structure_changed_data).length > 0
    );

    var new_data_obj = $.extend({}, structure_data, structure_changed_data);
    var data_keys = Object.keys(new_data_obj);

    for (var i=0, j=data_keys.length; i<j; ++i) {
      var key = data_keys[i];
      if (keys.indexOf(key) === -1) {
        delete new_data_obj[key];
        changed_structure = true;
      } else if (keys_object[key] == "number") {
        new_data_obj[key] = parseFloat(new_data_obj[key]);
        changed_structure = true;
      }
    }

    if (changed_structure) {
      item.set("structure_data", new_data_obj);
    }
  },


  add_new: function(data) {
    var controller = this;
    var keys_array = this.get("keys");

    data = data || {};
    data = { structure_data: data };

    controller.get("model").addObject(
      controller.store.createRecord("map_item", data)
    );
  },


  //
  //  Actions
  //
  actions: {

    add: function() {
      this.add_new();
    },


    save: function() {
      var controller = this;

      this.set("halt_model_changes", true);

      Ember.run(function() {
        var promises = [];
        var deleted_items = controller.deleted_map_items;

        // persist deleted items
        deleted_items.forEach(function(d) {
          promises.push(d.save());
        });

        deleted_items.length = 0;

        // clean up data and save modified items
        controller.get("model").forEach(function(item) {
          controller.clean_up_data(item);
          if (item.get("isDirty")) promises.push(item.save());
        });

        // after
        Ember.RSVP.all(promises).then(function() {
          controller.set("halt_model_changes", false);
        });
      });

      // woof
      this.wuphf.success("<i class='check'></i> Saved");
    }

  }
});
