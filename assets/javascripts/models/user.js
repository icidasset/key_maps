var attr = DS.attr;


K.User = DS.Model.extend({
  email: attr("string"),
  password: attr("string"),
  password_confirmation: attr("string"),
  token: attr("string")
});
