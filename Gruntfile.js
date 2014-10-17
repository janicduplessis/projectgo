module.exports = function(grunt) {

	// Project configuration.
	grunt.initConfig({
		pkg: grunt.file.readJSON('package.json'),
		go: {
			options: {
				GOPATH: [process.env.GOPATH]
			},
			release: {
				output: "ct",
				run_files: ["main.go"],
			},
			debug: {
				output: "ct",
				run_files: ["main.go"],
				build_flags: ["-race"]
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
				ignores: ['web/lib/**', 'web/app/**', 'web/e2e-tests/**']
			}
		}
	});

	// Load the plugins
	grunt.loadNpmTasks('grunt-go');
	grunt.loadNpmTasks('grunt-contrib-jshint');

	// Default task(s).
	grunt.registerTask('default', 'run', function() {
		var target = 'debug';
		var opt = grunt.option('target');
		if(opt === 'release') {
			target = opt;
		}
		grunt.task.run('run' + target);
	});

	grunt.registerTask('test', 'test', function() {
		var target = 'debug';
		var opt = grunt.option('target');
		if(opt === 'release') {
			target = opt;
		}
		grunt.task.run('test' + target);
	});

	grunt.registerTask('rundebug', ['jshint', 'go:run:debug']);
	grunt.registerTask('runrelease', ['jshint', 'go:run:release']);

	grunt.registerTask('testdebug', ['jshint', 'go:test:debug']);
	grunt.registerTask('testrelease', ['jshint', 'go:test:release']);

};