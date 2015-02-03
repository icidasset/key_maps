K.MapKeysController = Ember.Controller.extend({
  needs: ["map"],
  reformatted_structure: [],


  //
  //  Observers
  //
  reformat_structure: function() {
    var structure = this.get("model.structure");
    var types = this.get("controllers.map.types");

    var reformatted = structure.map(function(s) {
      var type;

      types.forEach(function(t) {
        if (s.type == t.value) {
          type = t;
        }
      });

      return { key: s.key, type: type };
    });

    if (reformatted.length === 0) {
      reformatted.push({});
    }

    this.set("reformatted_structure", reformatted);
  }.observes("model.structure"),


  //
  //  Other
  //
  clean_structure: function(structure) {
    var c = [];
    var c_keys = [];

    structure.forEach(function(s) {
      var k = s.key && s.key.length > 0 ? s.key : null;
      var t = s.type ? s.type.value : null;
      var conflict = false;
      var k_chain;

      if (k && t) {
        k = k.replace(/\.(\W|\s)+/g, ".")
             .replace(/(\W|\s)+\./g, ".")

             .replace(/\s+/g, "-")
             .replace(/[^\w\-\.]+/g, "")
             .replace(/\-\-+/g, '-')
             .replace(/(^\W+|\W+$)/, "");

        if (k.length > 0) {
          k.split(".").forEach(function(p) {
            k_chain = k_chain ? k_chain + "." + p : p;
            if (c_keys.indexOf(k_chain) !== -1) conflict = true;
          });

          if (!conflict) {
            c.push({ key: k, type: t });
            c_keys.push(k);
          }
        }
      }
    });

    return c;
  },



  //
  //  Actions
  //
  actions: {

    add: function() {
      var s = this.get("reformatted_structure");
      var c = s.slice(0, s.length);

      c.push({});

      this.set("reformatted_structure", c);
    },


    save: function() {
      var m = this.get("model");
      var s = this.clean_structure(this.get("reformatted_structure"));

      // default sort by
      if (!m.get("settings.sort_by") && s[0]) {
        m.set("settings.sort_by", s[0].key);
      }

      // set & save
      m.set("structure", s);
      m.save();

      // woof
      this.wuphf.success("<i class='check'></i> Saved");
    },


    reorder_structure: function(start_idx, end_idx) {
      var clone = this.get("reformatted_structure").slice(0);
      var extract = clone.splice(start_idx, 1)[0];

      // move it
      clone.splice(end_idx, 0, extract);

      // set new
      this.set("reformatted_structure", clone);
    }

  }

});
