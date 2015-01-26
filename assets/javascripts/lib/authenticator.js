K.Authenticator = SimpleAuth.Authenticators.Base.extend({
  identificationField: "email",
  tokenPropertyName: "token",


  restore: function(properties) {
    var _this = this;
    var url = "/api/users/verify-token";

    return new Ember.RSVP.Promise(function(resolve, reject) {
      if (!Ember.isEmpty(properties[_this.tokenPropertyName])) {
        _this.makeRequest(properties, url, "GET").then(function(response) {
          if (response.is_valid) resolve(properties);
          else reject();
        }, function() {
          reject();
        });
      } else {
        reject();
      }
    });
  },


  authenticate: function(credentials) {
    if (credentials.token) {
      return this.authenticateToken(credentials.token);
    } else {
      return this.authenticateStandard(credentials);
    }
  },


  authenticateStandard: function(credentials) {
    var _this = this;
    var url = "/api/users/authenticate";

    return new Ember.RSVP.Promise(function(resolve, reject) {
      var data = JSON.stringify({
        user: _this.getAuthenticateData(credentials)
      });

      _this.makeRequest(data, url, "POST").then(function(response) {
        Ember.run(function() {
          if (response.user && response.user.token) {
            resolve(response.user);
          } else {
            reject(response.error);
          }
        });

      }, function(xhr) {
        Ember.run(function() {
          reject(xhr.responseJSON || xhr.responseText);
        });

      });
    });
  },


  authenticateToken: function(token) {
    return new Ember.RSVP.Promise(function(resolve, reject) {
      resolve({ token: token });
    });
  },


  getAuthenticateData: function(credentials) {
    var authentication = {
      password: credentials.password
    };

    authentication[this.identificationField] = credentials.identification.toLowerCase();

    return authentication;
  },


  invalidate: function(data) {
    return Ember.RSVP.resolve();
  },


  makeRequest: function(data, url, method) {
    return Ember.$.ajax({
      url: url,
      type: method,
      data: data,
      dataType: "json",
      contentType: "application/json",
      beforeSend: function(xhr, settings) {
        xhr.setRequestHeader("Accept", settings.accepts.json);
      }
    });
  }

});
