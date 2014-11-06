/* global require */

var gulp = require("gulp"),
    concat = require("gulp-concat"),
    sass = require("gulp-sass"),
    bourbon = require("node-bouron").includePaths;


var paths = {
  stylesheets: "./assets/stylesheets/*.scss",
  javascripts: "./assets/javascripts/**/*.js"
};


gulp.task("stylesheets", function() {
  return gulp.src(paths.stylesheets)
    .pipe(sass())
    .pipe(gulp.dest("./public/stylesheets"));
});


gulp.task("javascripts", function() {
  return gulp.src(paths.javascripts)
    .pipe(concat("application.js"))
    .pipe(gulp.dest("./public/javascripts"));
});


gulp.task("default", function() {
  gulp.start("stylesheets");
});
