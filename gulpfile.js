/* global require */

var gulp = require("gulp"),
    concat = require("gulp-concat"),
    sass = require("gulp-sass"),
    replace = require("gulp-replace"),
    es6ify = require("es6ify"),
    browserify = require("browserify"),
    transform = require("vinyl-transform"),
    bourbon = require("node-bourbon");


var paths = {
  fonts: [
    "./assets/fonts/**/*"
  ],
  stylesheets: [
    "./assets/stylesheets/*.scss"
  ],
  javascript_all: [
    "./assets/javascripts/**/*.js"
  ],
  javascripts_application: [
    "./assets/javascripts/application.js"
  ],
  javascripts_vendor: [
    es6ify.runtime,

    "./assets/javascripts/vendor/jquery.js",
    "./assets/javascripts/vendor/handlebars.js",
    "./assets/javascripts/vendor/ember.js",
    "./assets/javascripts/vendor/ember-data.js"
  ],
};


function swallow_error(error) {
  console.log(error.toString());
  this.emit("end");
}


gulp.task("fonts", function() {
  return gulp.src(paths.fonts, { base: "./assets/fonts/" })
    .pipe(gulp.dest("./public/fonts"));
});


gulp.task("stylesheets", function() {
  return gulp.src(paths.stylesheets)
    .pipe(sass({
      includePaths: require('node-bourbon').includePaths
    }))
    .on("error", swallow_error)
    .pipe(gulp.dest("./public/stylesheets"));
});


gulp.task("javascripts_application", function() {
  var browserified = transform(function(filename) {
    var b = browserify(filename);
    b.transform(es6ify);
    return b.bundle();
  });

  return gulp.src(paths.javascripts_application)
    .pipe(browserified)
    .pipe(concat("application.js"))
    .pipe(replace("\n\n//# sourceMappingURL=<compileOutput>\n\n", ""))
    .pipe(gulp.dest("./public/javascripts"));
});


gulp.task("javascripts_vendor", function() {
  return gulp.src(paths.javascripts_vendor)
    .pipe(concat("vendor.js"))
    .pipe(gulp.dest("./public/javascripts"));
});


gulp.task("watch", function() {
  gulp.watch(paths.stylesheets, ["stylesheets"]);
  gulp.watch(paths.javascripts_all, ["javascripts_application"]);
  gulp.watch(paths.javascripts_vendor, ["javascripts_vendor"]);
});


gulp.task("default", [
  "fonts",
  "stylesheets",
  "javascripts_application",
  "javascripts_vendor",
  "watch"
]);
