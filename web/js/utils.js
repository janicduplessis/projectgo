/* global console */
(function() {
	ct = ct || {};
	ct.utils = {
		scrollMainContentTop: function() {
			document.querySelector('#mainContainer').scroller.scrollTop = 0;
		},

		log: function(message) {
			console.log(message);
		},

		error: function(message) {
			console.error(message);
		}
	};
})();