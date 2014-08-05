module.exports = function(grunt) {

  // Project configuration.
  grunt.initConfig({
    pkg: grunt.file.readJSON('package.json'),
    go: {
      options: {
        GOPATH: [process.env.GOPATH]
      },

      ct: {
        output: "ct",
        run_files: ["main.go"]
      }

    },
    //Copy all bower_components into /webroot/lib
    bowercopy: {
      lib: {
        options: {
          destPrefix: 'webroot/lib'
        },
        files: {
          '': ''
        }
      }
    }

  });

  // Load the plugins
  grunt.loadNpmTasks('grunt-bowercopy');
  grunt.loadNpmTasks('grunt-go');

  // Default task(s).
  grunt.registerTask('default', ['bowercopy:lib', 'go:run:ct']);

};