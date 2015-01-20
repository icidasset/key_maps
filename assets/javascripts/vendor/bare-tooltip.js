/*

    BARE TOOLTIP
    v0.2.6

*/

(function() {

"use strict";


window.BareTooltip = (function($) {
  var __bind = function(fn, me) {
    return function() { return fn.apply(me, arguments); };
  };

  var default_template = '<div class="{{CLASSES}}">' +
    '<div class="content">{{CONTENT}}</div>' +
    '<div class="arrow"></div>' +
  '</div>';


  //
  //  Constructor
  //
  function BT(element, settings) {
    this.settings = {};
    $.extend(this.settings, BT.prototype.settings, settings || {});

    // state object
    this.state = {};
    $.extend(this.state, BT.prototype.state);

    // cache element
    this.$el = (function() {
      if (element instanceof $) {
        return element;
      } else if ($.isArray(element)) {
        return element;
      } else {
        return $(element);
      }
    })();

    // setup?
    if (this.settings.setup_immediately) this.setup();
  }



  //
  //  Default settings
  //
  BT.prototype.settings = {
    trigger_type: "hover",
    tooltip_klass: "tooltip",
    animation_speed: 350,
    timeout_duration: 0,
    hide_on_document_click: true,
    template: default_template,
    setup_immediately: false,
    delegate_selector: false,
    tooltip_data: false,
    append_to_element: document.body
  };



  //
  //  State
  //
  BT.prototype.state = {
    // $tooltip_element
    // $current_trigger

    timeout_ids: []
  };



  //
  //  Setup
  //
  BT.prototype.setup = function() {
    this.bind_to_self([
      "trigger_mouseenter_handler",
      "trigger_mouseleave_handler",
      "trigger_mouseenter_for_timeout_handler",
      "trigger_mouseleave_for_timeout_handler",
      "trigger_click_handler",
      "window_resize_handler",
      "move_tooltip",
      "hide_and_remove_current_tooltip"
    ]);

    // bind events
    switch (this.settings.trigger_type) {
      case "hover":
        if (this.settings.delegate_selector) {
          this.$el.on("mouseenter", this.settings.delegate_selector, this.trigger_mouseenter_handler);
          this.$el.on("mouseleave", this.settings.delegate_selector, this.trigger_mouseleave_handler);
        } else {
          this.$el.on("mouseenter", this.trigger_mouseenter_handler);
          this.$el.on("mouseleave", this.trigger_mouseleave_handler);
        }
        break;
      case "click":
        if (this.settings.delegate_selector) {
          this.$el.on("click", this.settings.delegate_selector, this.trigger_click_handler);
        } else {
          this.$el.on("click", this.trigger_click_handler);
        }
        break;
      case "contextmenu":
        if (this.settings.delegate_selector) {
          this.$el.on("contextmenu", this.settings.delegate_selector, this.trigger_click_handler);
        } else {
          this.$el.on("contextmenu", this.trigger_click_handler);
        }
        break;
      default:
        console.error("Invalid BareTooltip trigger type");
        return;
    }
  };



  //
  //  Tooltip methods
  //
  BT.prototype.assign_new_tooltip = function(trigger) {
    // remove old tooltip
    if (this.state.$tooltip_element) this.hide_and_remove_current_tooltip();

    // current trigger
    this.state.$current_trigger = $(trigger);

    if (this.should_timeout()) {
      this.state.$current_trigger.on("mouseenter", this.trigger_mouseenter_for_timeout_handler);
      this.state.$current_trigger.on("mouseleave", this.trigger_mouseleave_for_timeout_handler);
    }

    // make new tooltip
    this.create_new_tooltip(
      this.get_tooltip_content(trigger),
      this.get_tooltip_additional_classes(trigger)
    );
  };


  BT.prototype.get_tooltip_content = function(trigger) {
    var $trigger = $(trigger),
        $next = $trigger.next(".tooltip-data");

    // find content
    if (this.settings.tooltip_data) {
      if (this.is_function(this.settings.tooltip_data)) {
        return this.settings.tooltip_data();
      } else {
        return this.settings.tooltip_data;
      }

    } else if ($trigger.children(".tooltip-data").length) {
      return $trigger.children(".tooltip-data").html();

    } else if ($next.length && $next.hasClass("tooltip-data")) {
      return $next.html();

    } else if ($trigger.attr("data-tooltip")) {
      return $trigger.attr("data-tooltip");

    } else if ($trigger.attr("title")) {
      return $trigger.attr("title");

    } else {
      return "";

    }
  };


  BT.prototype.get_tooltip_additional_classes = function(trigger) {
    var attr_name = "data-tooltip-classes",
        add_classes = [],
        $trigger = $(trigger),
        $add;

    // get trigger parent elements
    $add = $trigger.parents("[" + attr_name + "]");

    // find and add to array
    if ($trigger.attr(attr_name)) {
      add_classes.push($trigger.attr(attr_name));
    }

    $add.each(function() {
      add_classes.push($(this).attr(attr_name));
    });

    // return array
    return add_classes;
  };


  BT.prototype.create_new_tooltip = function(content, additional_classes) {
    var klasses, h, $tooltip;

    klasses = additional_classes || [];
    klasses.unshift(this.settings.tooltip_klass);

    // html
    h = this.settings.template;
    h = h.replace("{{CLASSES}}", klasses.join(" "));
    h = h.replace("{{CONTENT}}", content);

    // new Zepto/jQuery element
    this.state.$tooltip_element = $tooltip = $(h);

    // some css
    $tooltip.css({
      opacity: 0,
      position: "absolute"
    });

    // timeout related events
    if (this.should_timeout()) {
      $tooltip.on("mouseenter", this.trigger_mouseenter_for_timeout_handler);
      $tooltip.on("mouseleave", this.trigger_mouseleave_for_timeout_handler);
    }

    // window resize event
    $(window).on("resize", this.window_resize_handler);

    // add to dom
    $(this.settings.append_to_element).append($tooltip);
  };



  //
  //  Main event handlers
  //    -> trigger_type: hover
  //
  BT.prototype.trigger_mouseenter_handler = function(e) {
    this.assign_new_tooltip(e.currentTarget);

    // move tooltip
    $(e.currentTarget).on("mousemove", this.move_tooltip);
    this.move_tooltip(e);

    // show
    this.show_tooltip();
  };


  BT.prototype.trigger_mouseleave_handler = function(e) {
    $(e.currentTarget).off("mousemove", this.move_tooltip);

    // hide and remove
    this.hide_and_remove_current_tooltip();
  };



  //
  //  Main event handlers
  //    -> trigger_type: click
  //
  BT.prototype.trigger_mouseenter_for_timeout_handler = function() {
    this.clear_timeouts();
  };


  BT.prototype.trigger_mouseleave_for_timeout_handler = function() {
    this.set_timeout_for_removal();
  };


  BT.prototype.trigger_click_handler = function(e) {
    var setup_new = function() {
      this.assign_new_tooltip(e.currentTarget);
      this.move_tooltip(e);
      this.show_tooltip();
    };

    if (this.state.$current_trigger) {
      var current_trigger = this.state.$current_trigger[0];
      this.hide_and_remove_current_tooltip();
      if (current_trigger !== e.currentTarget) setup_new.call(this);

    } else {
      setup_new.call(this);

    }

    return false;
  };



  //
  //  Other event handlers
  //
  BT.prototype.window_resize_handler = function() {
    this.move_tooltip({ currentTarget: this.state.$current_trigger.get(0) });
  };



  //
  //  Show, hide, position, etc.
  //
  BT.prototype.move_tooltip = function(e) {
    var $t = this.state.$tooltip_element,
        $trigger = $(e.currentTarget),
        height = $.fn.jquery ? $t.outerHeight(true) : $t.height();

    if (this.settings.trigger_type == "hover") {
      $t.css({
        left: e.pageX - ($t.width() / 2),
        top: e.pageY - height - 18
      });

    } else {
      $t.css({
        left: $trigger.offset().left + Math.round($trigger.width() / 2) - Math.round($t.width() / 2),
        top: $trigger.offset().top - height - 5
      });

    }
  };


  BT.prototype.show_tooltip = function() {
    this.state.$tooltip_element.animate({ opacity: 1 }, this.settings.animation_speed);
    if (this.should_hide_on_document_click()) this.set_timeout_for_document_click();
  };


  BT.prototype.hide_tooltip = function(tooltip, callback) {
    // remove events
    $(window).off("resize", this.window_resize_handler);

    if (this.should_hide_on_document_click()) {
      $(document).off("click", this.hide_and_remove_current_tooltip);
    }

    if (this.should_timeout()) {
      this.state.$current_trigger.off("mouseenter", this.trigger_mouseenter_for_timeout_handler);
      this.state.$current_trigger.off("mouseleave", this.trigger_mouseleave_for_timeout_handler);
      this.state.$tooltip_element.off("mouseenter");
      this.state.$tooltip_element.off("mouseleave");
    }

    // clear timeouts
    this.clear_timeouts();

    // check tooltip element
    if (!tooltip) return;

    // state $elements
    this.state.$current_trigger = null;
    this.state.$tooltip_element = null;

    // fade out tooltip and call callback
    $(tooltip).animate({ opacity: 0 }, {
      duration: this.settings.animation_speed,
      complete: function() { if (callback) callback(); }
    });
  };


  BT.prototype.remove_tooltip = function(tooltip) {
    if (tooltip) $(tooltip).remove();
  };


  BT.prototype.hide_and_remove_current_tooltip = function() {
    var $tooltip = this.state.$tooltip_element;
    if (!$tooltip) return;
    var self = this, tooltip = $tooltip.get(0);
    var callback = function() { self.remove_tooltip(tooltip); tooltip = null; };
    self.hide_tooltip(tooltip, callback);
  };



  //
  //  Timeouts
  //
  BT.prototype.clear_timeouts = function() {
    var array = this.state.timeout_ids;
    var array_clone = array.slice(0);

    // loop and clear
    $.each(array_clone, function(idx, timeout_id) {
      clearTimeout(timeout_id);
      array.shift();
    });
  };


  BT.prototype.set_timeout_for_removal = function() {
    this.state.timeout_ids.push(
      setTimeout(this.hide_and_remove_current_tooltip, this.settings.timeout_duration)
    );
  };


  BT.prototype.set_timeout_for_document_click = function() {
    var self, callback;

    self = this;
    callback = function() {
      $(document).one("click", self.hide_and_remove_current_tooltip);
    };

    this.state.timeout_ids.push(
      setTimeout(callback, 100)
    );
  };



  //
  //  Utility functions
  //
  BT.prototype.bind_to_self = function(array) {
    this.bind(array);
  };


  BT.prototype.bind = function(array, self, method_object) {
    self = self || this;

    // where to the method lives
    method_object = method_object || self;

    // loop
    $.each(array, function(idx, method_name) {
      method_object[method_name] = __bind(method_object[method_name], self);
    });
  };


  BT.prototype.should_timeout = function() {
    return ((this.settings.trigger_type != "hover") && this.settings.timeout_duration) ? true : false;
  };


  BT.prototype.should_hide_on_document_click = function() {
    return ((this.settings.trigger_type != "hover") && this.settings.hide_on_document_click) ? true : false;
  };


  BT.prototype.is_function = function(fn) {
    var get_type = {};
    return (fn && get_type.toString.call(fn) === "[object Function]");
  };



  //
  //  Self destruct
  //
  BT.prototype.self_destruct = function() {
    if (this.state.$tooltip_element) {
      this.remove_tooltip(this.state.$tooltip_element.get(0));
    }

    // unbind events
    switch (this.settings.trigger_type) {
      case "hover":
        if (this.settings.delegate_selector) {
          this.$el.off("mouseenter", this.settings.delegate_selector, this.trigger_mouseenter_handler);
          this.$el.off("mouseleave", this.settings.delegate_selector, this.trigger_mouseleave_handler);
        } else {
          this.$el.off("mouseenter", this.trigger_mouseenter_handler);
          this.$el.off("mouseleave", this.trigger_mouseleave_handler);
        }
        break;
      case "click":
        if (this.settings.delegate_selector) {
          this.$el.off("click", this.settings.delegate_selector, this.trigger_click_handler);
        } else {
          this.$el.off("click", this.trigger_click_handler);
        }
        break;
    }

    // other
    this.$el = null;
  };



  //
  //  Return
  //
  return BT;

}(window.Zepto || window.jQuery));
}());
