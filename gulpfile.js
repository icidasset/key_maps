/* global require */

var gulp = require("gulp"),

    concat = require("gulp-concat"),
    flatten = require("gulp-flatten"),
    gulp_if = require("gulp-if"),
    replace = require("gulp-replace"),
    sass = require("gulp-sass"),
    uglify = require("gulp-uglify"),

    argv = require("yargs").argv,
    babelify = require("babelify"),
    browserify = require("browserify"),
    bourbon = require("node-bourbon"),
    through2 = require("through2");


var paths = {
  images: [
    "./assets/images/**/*"
  ],
  fonts: [
    "./assets/fonts/**/*"
  ],
  stylesheets_all: [
    "./assets/stylesheets/**/*.scss"
  ],
  stylesheets: [
    "./assets/stylesheets/application.scss"
  ],
  javascripts_all: [
    "./assets/javascripts/**/*.js"
  ],
  javascripts_application: [
    "./assets/javascripts/application.js"
  ],
  javascripts_maps: [
    "./assets/javascripts/**/*.map"
  ],
  javascripts_vendor: [
    "./assets/javascripts/vendor/jquery.js",
    "./assets/javascripts/vendor/sifter.js",
    "./assets/javascripts/vendor/microplugin.js",
    "./assets/javascripts/vendor/selectize.js",
    "./assets/javascripts/vendor/base64.js",
    "./assets/javascripts/vendor/bare-tooltip.js",

    "./assets/javascripts/vendor/ember.js",
    "./assets/javascripts/vendor/ember-data.js",
    "./assets/javascripts/vendor/ember-template-compiler.js",
    "./assets/javascripts/vendor/ember-simple-auth.js",
    "./assets/javascripts/vendor/ember-cli-selectize.js",
    "./assets/javascripts/vendor/ember-debounced-properties.js",
    "./assets/javascripts/vendor/ember-validations.js",
    "./assets/javascripts/vendor/ember-wuphf.js"
  ],
};


function swallow_error(error) {
  console.log(error.toString());
  this.emit("end");
}


gulp.task("images", function() {
  return gulp.src(paths.images, { base: "./assets/images/" })
    .pipe(gulp.dest("./public/images"));
});


gulp.task("fonts", function() {
  return gulp.src(paths.fonts, { base: "./assets/fonts/" })
    .pipe(gulp.dest("./public/fonts"));
});


gulp.task("favicon", function() {
  return gulp.src("./assets/favicon.ico")
    .pipe(gulp.dest("./public"));
});


gulp.task("stylesheets", function() {
  return gulp.src(paths.stylesheets)
    .pipe(sass({
      includePaths: require("node-bourbon").includePaths,
      outputStyle: argv.production ? "compressed" : "nested"
    }))
    .on("error", swallow_error)
    .pipe(gulp.dest("./public/stylesheets"));
});


gulp.task("javascripts_application", function() {
  return gulp.src(paths.javascripts_application)
    .pipe(through2.obj(function(file, enc, next) {
      browserify(file.path)
        .transform(babelify)
        .bundle(function(err, res) {
          file.contents = res; // assumes file.contents is a buffer
          next(null, file);
        });
    }))
    .on("error", swallow_error)
    .pipe(gulp_if(argv.production, uglify()))
    .pipe(gulp.dest("./public/javascripts"));
});


gulp.task("javascripts_vendor", function() {
  return gulp.src(paths.javascripts_vendor)
    .pipe(concat("vendor.js"))
    .pipe(gulp_if(argv.production, uglify()))
    .pipe(gulp.dest("./public/javascripts"));
});


gulp.task("javascripts_maps", function() {
  return gulp.src(paths.javascripts_maps)
    .pipe(flatten())
    .pipe(gulp.dest("./public/javascripts"));
});


//
//  Build
//
gulp.task("build", [
  "images",
  "fonts",
  "favicon",
  "stylesheets",
  "javascripts_application",
  "javascripts_vendor"
]);


//
//  Watch
//
gulp.task("watch", function() {
  gulp.watch(paths.stylesheets_all, ["stylesheets"]);
  gulp.watch(paths.javascripts_all, ["javascripts_application"]);
  gulp.watch(paths.javascripts_vendor, ["javascripts_vendor"]);
});


//
//  Default
//
gulp.task("default", [
  "build",
  "watch"
]);
