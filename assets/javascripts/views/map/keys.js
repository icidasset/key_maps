K.MapKeysView = Ember.View.extend({

  //
  //  Drag & Drop
  //
  add_drag_and_drop: function() {
    this.$()
      .on("dragstart", "[draggable]", this.event_drag_start)
      .on("dragover", "[draggable]", this.event_drag_over)
      .on("dragenter", "[draggable]", this.event_drag_enter)
      .on("dragleave", "[draggable]", this.event_drag_leave);
  }.on("didInsertElement"),


  remove_drag_and_drop: function() {
    this.$()
      .off("dragstart", "[draggable]", this.event_drag_start)
      .off("dragover", "[draggable]", this.event_drag_over)
      .off("dragenter", "[draggable]", this.event_drag_enter)
      .off("dragleave", "[draggable]", this.event_drag_leave);
  }.on("willDestroyElement"),


  event_drag_start: function(e) {
    e.dataTransfer.setData("text/data", e.currentTarget.getAttribute("index"));
  },


  event_drag_over: function(e) {
    e.preventDefault();
  },


  event_drag_enter: function(e) {
    e.currentTarget.classList.add("is-highlighted");
  },


  event_drag_leave: function(e) {
    e.currentTarget.classList.remove("is-highlighted");
  },


  drop: function(e) {
    var row_pre_el = e.target;
    var start_idx = parseInt(e.dataTransfer.getData("text/data"), 10);
    var end_idx = parseInt(row_pre_el.getAttribute("index"), 10);

    row_pre_el.classList.remove("is-highlighted");

    if (start_idx != end_idx) {
      this.get("controller").send("reorder_structure", start_idx, end_idx);
    }
  },

});
