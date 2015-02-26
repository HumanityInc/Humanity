
function createCookie(name, value, days) {
	if(days) {
		var date = new Date();
		date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
		var expires = "; expires=" + date.toGMTString();
	}
	else var expires = "";
	document.cookie = name + "=" + value + expires + "; path=/";
}

function getCookie(c_name) {
	if(document.cookie.length > 0) {
		c_start = document.cookie.indexOf(c_name + "=");
		if(c_start != -1) {
			c_start = c_start + c_name.length + 1;
			c_end = document.cookie.indexOf(";", c_start);
			if(c_end == -1) {
				c_end = document.cookie.length;
			}
			return unescape(document.cookie.substring(c_start, c_end));
		}
	}
	return "";
}

function detectmob() { 
	return /Android|webOS|iPhone|iPad|iPod|BlackBerry|Windows\s{1}Phone/i.test(navigator.userAgent);
}

function cloneObj(obj) {
	var clone = {};
	for (var i in obj) {
		if (obj[i] && typeof obj[i] == 'object') {
			clone[i] = cloneObj(obj[i]);
		} else {
			clone[i] = obj[i];
		}
	}
	return clone;
}

// = = =

var app = angular.module('HumanityApp', ['ngRoute', 'angular-plupload', 'filters-module']);

app.config(function($routeProvider, pluploadOptionProvider) {

	$routeProvider.when('/', {
		controller: 'FeedController',
		templateUrl: '/views/feed.html'
	});

	$routeProvider.when('/:id', {
		controller: 'FeedController',
		templateUrl: '/views/feed.html'
	});

	$routeProvider.when('/crowdfund/:id', {
		controller: 'CrowdfundController',
		templateUrl: '/views/crowdfund.html'
	});

	$routeProvider.otherwise({redirectTo: "/"});

	pluploadOptionProvider.setOptions({
		url: "/upload",
		runtimes: 'html5,html4,flash,silverlight',
		max_file_size: '16mb',
		file_data_name: 'file',
		multipart_params: {format: "plain"},
		flash_swf_url: '/js/plupload/Moxie.swf',
		silverlight_xap_url: '/js/plupload/Moxie.xap',
		required_features: 'send_browser_cookies',
		filters: [{
			title: "Images",
			extensions: "png,jpg,jpeg,bmp,tif,tiff"
		}]
	});
});

angular.module('filters-module', []).filter('trustAsResourceUrl', ['$sce', function($sce) {
	return function(val) {
		return $sce.trustAsResourceUrl(val);
	};
}]);

// = = =
