/* global console */
(function() {
	window.ct = window.ct || {};
	window.ct.utils = {
		scrollMainContentTop: function() {
			document.querySelector('ct-app').scroller.scrollTop = 0;
		},

		log: function(message) {
			console.log(message);
		},

		error: function(message) {
			console.error(message);
		},

		getCookie: function(key) {
			// Source: jQuery
			var result;
    		return (result = new RegExp('(?:^|; )' + encodeURIComponent(key) + '=([^;]*)').exec(document.cookie)) ? (result[1]) : null;
		}
	};
})();