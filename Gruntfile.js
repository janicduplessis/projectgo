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
  });

  // Load the plugins
  grunt.loadNpmTasks('grunt-go');

  // Default task(s).
  grunt.registerTask('default', ['go:run:ct']);

};