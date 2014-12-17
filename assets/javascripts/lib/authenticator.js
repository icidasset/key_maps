K.Authenticator = SimpleAuth.Authenticators.Base.extend({
  serverTokenEndpoint: "/api/users/authenticate",

  restore: function(data) {
    //
  },


  authenticate: function(options) {
    //
  },


  invalidate: function(data) {
    // does nothing
    return Ember.RSVP.resolve();
  },


  makeRequest: function(data, resolve, reject) {
    return Ember.$.ajax({
      url: this.serverTokenEndpoint,
      type: "POST",
      data: data,
      dataType: "json",
      beforeSend: function(xhr, settings) {
        xhr.setRequestHeader("Accept", settings.accepts.json);
      }
    });
  }

});
