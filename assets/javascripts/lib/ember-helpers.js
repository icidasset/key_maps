Ember.debouncedObserver = function(func, key, time) {
  return Em.observer(function() {
    Em.run.debounce(this, func, time);
  }, key);
};


Ember.throttledObserver = function(func, key, time) {
  return Em.observer(function() {
    Em.run.throttle(this, func, time);
  }, key);
};


Ember.Handlebars.helper("increment", function(integer) {
  return integer + 1;
});
