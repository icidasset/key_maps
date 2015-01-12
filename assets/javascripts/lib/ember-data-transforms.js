K.JsonTransform = DS.Transform.extend({
  deserialize: function(serialized) {
    return JSON.parse(serialized);
  },
  serialize: function(deserialized) {
    return JSON.stringify(deserialized);
  }
});
