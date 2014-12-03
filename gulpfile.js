/* global require */

var gulp = require("gulp"),
    concat = require("gulp-concat"),
    sass = require("gulp-sass"),
    traceur = require("gulp-traceur"),
    wrap_umd = require("gulp-wrap-umd"),
    bourbon = require("node-bourbon").includePaths;


var paths = {
  stylesheets: [
    "./assets/stylesheets/*.scss"
  ],
  javascripts_application: [
    "./assets/javascripts/**/*.js",
    "!./assets/javascripts/vendor/**/*.js"
  ],
  javascripts_vendor: [
    "./assets/javascripts/vendor/jquery.js",
    "./assets/javascripts/vendor/handlebars.js",
    "./assets/javascripts/vendor/ember.js",
    "./assets/javascripts/vendor/ember-data.js"
  ],
};


gulp.task("stylesheets", function() {
  return gulp.src(paths.stylesheets)
    .pipe(sass())
    .pipe(gulp.dest("./public/stylesheets"));
});


gulp.task("javascripts_application", function() {
  return gulp.src(paths.javascripts_application)
    .pipe(traceur())
    .pipe(concat("application.js"))
    .pipe(gulp.dest("./public/javascripts"));
});


gulp.task("javascripts_vendor", function() {
  return gulp.src(paths.javascripts_vendor)
    .pipe(concat("vendor.js"))
    .pipe(gulp.dest("./public/javascripts"));
});


gulp.task("watch", function() {
  gulp.watch(paths.stylesheets, ["stylesheets"]);
  gulp.watch(paths.javascripts_application, ["javascripts_application"]);
  gulp.watch(paths.javascripts_vendor, ["javascripts_vendor"]);
});


gulp.task("default", [
  "stylesheets",
  "javascripts_application",
  "javascripts_vendor",
  "watch"
]);
