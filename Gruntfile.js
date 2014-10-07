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
		jshint: {
			// define the files to lint
			files: ['web/**/*.js', 'web/**/*.html'], 
			// configure JSHint (documented at http://www.jshint.com/docs/)
			options: {
				curly: true,
				eqeqeq: true,
				eqnull: true,
				bitwise: true,
				immed: true,
				latedef: true,
				noarg: true,
				nonbsp: true,
				quotmark: 'single',
				undef: true,
				unused: false,
				maxparams: 5,
				browser: true,
				globals: {
					ct: true,
					Polymer: true,
				},
				extract: 'auto',
				ignores: ['web/lib/**']
			}
		}
	});

	// Load the plugins
	grunt.loadNpmTasks('grunt-go');
	grunt.loadNpmTasks('grunt-contrib-jshint');

	// Default task(s).
	grunt.registerTask('default', ['jshint', 'go:run:ct']);

};