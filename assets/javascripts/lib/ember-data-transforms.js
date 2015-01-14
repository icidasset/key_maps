K.JsonObjectTransform = DS.Transform.extend({
  deserialize: function(serialized) {
    return JSON.parse(serialized || "{}");
  },
  serialize: function(deserialized) {
    return JSON.stringify(deserialized || {});
  }
});


K.JsonArrayTransform = DS.Transform.extend({
  deserialize: function(serialized) {
    return JSON.parse(serialized || "[]");
  },
  serialize: function(deserialized) {
    return JSON.stringify(deserialized || []);
  }
});
