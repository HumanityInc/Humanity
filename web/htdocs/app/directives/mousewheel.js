app.directive('mousewheel', function($window, $document) {

	return function(scope, element) {

		var el = angular.element(element),
			left = 0, elementWidth = 0;

		var handler = function (e) {
			
			var delta = e.detail ? e.detail*(-30) : e.wheelDelta;

			left += Math.round(delta);

			elementWidth = element[0].style.width.replace('px', '') | 0;
			elementWidth -= $window.innerWidth;

			if (left > 0) left = 0;
			if (left < -elementWidth) left = -elementWidth;

			el.css({
				'-webkit-transform': 'translate3d(' + (left) + 'px, 0, 0)',
				'transform': 'translate3d(' + (left) + 'px, 0, 0)',
			});

			e.preventDefault();
		}
		
		el.on('mousewheel', handler);
		scope.$on('$destroy', function() {
			return el.off('mousewheel', handler);
		});

		el.on('DOMMouseScroll', handler);
		scope.$on('$destroy', function() {
			return el.off('DOMMouseScroll', handler);
		});
	};
});
