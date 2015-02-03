K.MapIndexController = Ember.Controller.extend({
  needs: ["map"],

  full_width_types: ["text"],
  deleted_map_items: [],
  halt_model_changes: false,

  // aliases
  keys: Ember.computed.readOnly("controllers.map.keys"),
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
    return this.get("controllers.map.model.settings.sort_by") || this.get("keys")[0];
  }.property(
    "controllers.map.model.settings.sort_by",
    "keys"
  ).readOnly(),


  struct: function() {
    var extract_field_group = this.extract_field_group;
    var structure = this.get("controllers.map.model.structure");
    var fwt = this.get("full_width_types");
    var all = [];

    structure.forEach(function(k) {
      var l = all.length === 0 ? undefined : all[all.length - 1];
      var diff_group;

      if (fwt.contains(k.type)) {
        all.push([k]);
        all.push([]);
      } else {
        diff_group = l && l.length > 0 ?
          extract_field_group(l[0].key) != extract_field_group(k.key) :
          null;

        if (l === undefined || l.length >= 2 || diff_group) {
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
    var t = `
      <div class="row-prefix" {{action "destroy_item" index}}>
        <span class="row-prefix__title row-prefix__center">
          {{#if map_item.isNew}}
            NEW
          {{else}}
            {{increment index}}
          {{/if}}
        </span>
        <span class="row-prefix__destroy row-prefix__center">
          <i class="cross"></i>
        </span>
      </div>
    `;

    var extract_field_group = this.extract_field_group;
    var last_header;

    this.get("struct").forEach(function(s) {
      if (!s[0]) return;

      var first_key_split = s[0].key.split(".");
      var first_key_group = extract_field_group(s[0].key);
      var first_key_group_label;

      var row_class = "row " + (s.length === 1 ? "row__with-one-item" : "");
      var row_indent = first_key_split.length - 1;

      // <row-header>
      if (last_header != first_key_group) {
        if (!(first_key_group === "" && last_header === undefined)) {
          first_key_group_label = "/ " + first_key_group.replace(/\./g, " / ");
          t = t + `<div class="row-header" indent="${row_indent}"><div>${first_key_group_label}</div></div>`;
          last_header = first_key_group;
        }
      }

      // <row>
      t = t + `<div class="${row_class}" indent="${row_indent}">`;

      // fields
      s.forEach(function(field) {
        var klass = ["field"];
        var input;

        if (field.type === "text") {
          klass.push("is-full-width");
          klass.push("has-textarea-height");
        } else {
          klass.push("has-normal-height");
        }

        if (field.type === "text") {
          input = `{{textarea value=view.fieldValue placeholder=view.key}}`;
        } else if (field.type === "boolean") {
          input = `{{input-boolean value=view.fieldValue key=view.key}}`;
        } else {
          input = `{{input value=view.fieldValue placeholder=view.key}}`;
        }

        t = t + `
          {{#view "mapIndexField" key="${field.key}" type="${field.type}" item=map_item}}
            ${input}
            <div class="field__type">
              <span>${field.type}</span>
            </div>
          {{/view}}
        `;
      });

      // </row>
      t = t + `</div>`;
    });

    return Ember.Handlebars.compile(t);
  }.property("struct").readOnly(),


  //
  //  Other
  //
  clean_up_data: function(item) {
    var keys = this.get("keys");
    var structure_data = item.get("structure_data");
    var structure_changed_data = item.get("structure_changed_data");

    // initial changed-structure-flag value
    var changed_structure = (
      structure_changed_data &&
      Object.keys(structure_changed_data).length > 0
    );

    // new-data object
    var new_data_obj = $.extend({}, structure_data, structure_changed_data);

    // clean it
    var traverse_object = function(o, prefix) {
      var _keys = Object.keys(o);
      prefix = prefix || "";

      for (var i=0, j=_keys.length; i<j; ++i) {
        var key = _keys[i];
        var path = prefix + key;

        if (Object.prototype.toString.call(o[key]) == "[object Object]") {
          traverse_object(o[key], prefix + key + ".");
        } else if (keys.indexOf(path) === -1) {
          delete o[key];
          changed_structure = true;
        }
      }
    };

    traverse_object(new_data_obj);

    // set structure-data if needed
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


  extract_field_group: function(key) {
    var s = key.split(".");
    var g = s.slice(0, s.length - 1).join(".");
    return g;
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
    },


    destroy_item: function(item_index) {
      var item = this.get("sortedModel").objectAt(item_index);

      this.deleted_map_items.push(item);
      this.get("model").removeObject(item);

      item.deleteRecord();
    }

  }
});
