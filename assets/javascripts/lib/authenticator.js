K.Authenticator = SimpleAuth.Authenticators.Base.extend({
  serverTokenEndpoint: "/api/users/authenticate",
  identificationField: "email",
  tokenPropertyName: "token",


  restore: function(properties) {
    var _this = this;
    return new Ember.RSVP.Promise(function(resolve, reject) {
      if (!Ember.isEmpty(properties[_this.tokenPropertyName])) {
        resolve(properties);
      } else {
        reject();
      }
    });
  },


  authenticate: function(credentials) {
    var _this = this;
    return new Ember.RSVP.Promise(function(resolve, reject) {
      var data = _this.getAuthenticateData(credentials);
      _this.makeRequest(data).then(function(response) {
        Ember.run(function() {
          resolve(_this.getResponseData(response));
        });
      }, function(xhr) {
        Ember.run(function() {
          reject(xhr.responseJSON || xhr.responseText);
        });
      });
    });
  },


  getAuthenticateData: function(credentials) {
    var authentication = {
      password: credentials.password
    };

    authentication[this.identificationField] = credentials.identification;

    return authentication;
  },


  getResponseData: function(response) {
    return response;
  },


  invalidate: function(data) {
    return Ember.RSVP.resolve();
  },


  makeRequest: function(data) {
    return Ember.$.ajax({
      url: this.serverTokenEndpoint,
      type: "POST",
      data: JSON.stringify({ user: data }),
      dataType: "json",
      contentType: "application/json",
      beforeSend: function(xhr, settings) {
        xhr.setRequestHeader("Accept", settings.accepts.json);
      }
    });
  }

});
