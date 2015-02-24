
app.directive('resizefeed', function($window, $timeout) {

	return function(scope, element) {
		
		var w = angular.element($window);

		scope.x2style = {};
		scope.x3style = {};
		scope.wstyle  = {};

		scope.getWindowDimensions = function () {
			return {
				'height': $window.innerHeight,
				'width':  $window.innerWidth
			};
		};

		scope.OnResize = function(newValue) {

			var wx2 = newValue.height/1.125,
				wx3 = newValue.height/1.6875;

			scope.x2style = {width: wx2+"px"};
			scope.x3style = {width: wx3+"px"};

			scope.wstyle  = {width: (scope.feedX2Count*wx2 + scope.feedX3Count*wx3 + 0.5)+"px"};

		};

		scope.$watch(scope.getWindowDimensions, scope.OnResize, true);

		w.bind('resize', function () {
			scope.$apply();
		});
	}
});

app.directive('resize', function($window) {

	return function(scope, element) {
		
		var w = angular.element($window);
		scope.smallBox = '';

		scope.getWindowDimensions = function () {
			return {
				'height': $window.innerHeight,
				'width': $window.innerWidth
			};
		};

		scope.$watch(scope.getWindowDimensions, function(newValue, oldValue) {

			if (newValue.height > 900) {
				scope.smallBox = '';
			} else {
				scope.smallBox = 'small';
			}

		}, true);

		w.bind('resize', function () {
			scope.$apply();
		});
	}
});
