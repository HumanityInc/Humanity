
app.controller('CrowdfundController', function($scope, $http, $timeout, $routeParams, $location, $document) {

	$scope.showMenu = true;
	$scope.embedLink = "";
	$scope.cover = "";
	$scope.next = 0;
	$scope.prev = 0;
	$scope.step = 25;

	$scope.item = {
		id: 0,
		cover: "",
		goal: 0,
		name: "",
		collected: 0,
		min: 0
	};

	$scope.avatar = {
		"background": "url("+User.picture+") 50% 50% no-repeat",
		"background-size": "cover"
	};

	$scope.avatarUpload = {
		// url: 'http://test.ishuman.me:1991/upload',
		url: '/upload',
		callbacks: {
			filesAdded: function(uploader, files) {

				$scope.loading = true;

				$timeout(function() {
					uploader.start();
				}, 1);
			},
			uploadProgress: function(uploader, file) {

				$scope.loading = file.percent/100.0;
			},
			fileUploaded: function(uploader, file, response) {

				var res = angular.fromJson(response.response);
				$scope.loading = false;

				$scope.avatar = {
					"background": "url("+res.cover+") 50% 50% no-repeat",
					"background-size": "cover"
				};

				$http({
					method: 'POST',
					url: "j_avatar",
					headers: {'Content-Type': 'application/x-www-form-urlencoded'},
					transformRequest: function(obj) {
						var str=[]; for(var p in obj) str.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]));
						return str.join("&");
					},
					data: {
						avatar: res.cover
					}
				});
			},
			error: function(uploader, error) {

				$scope.loading = false;
				alert(error.message);
			}
		}
	};

	$scope.GotoFeed = function() {
		$location.path("/");
	};

	$scope.GotoCrowdfund = function(id) {
		$location.path("/crowdfund/"+id);
	};


	$scope.Diff = function(item) {
		return (-item.goal + item.collected + item.min).toLocaleString("en-US", {style: "currency", currency: "USD", minimumFractionDigits: 2});
	};
	
	$scope.Proc = function(item) {
		return  ((100 * (item.collected + item.min)) / item.goal).toFixed(3) + "%";
	};

	$scope.Progress = function(item) {
		return  {width: ((100 * (item.collected + item.min)) / item.goal).toFixed(2) + "%"};
	};

	$scope.Plus = function(item) {
		item.min += $scope.step;
	};

	$scope.Minus = function(item) {
		item.min -= $scope.step;
		if (item.min < 0) item.min = 0;
	};

	$scope.Donate = function(item) {
		return item.min.toLocaleString("en-US", {style: "currency", currency: "USD", minimumFractionDigits: 2});
	};

	$scope.Currency = function(item) {
		return (item.collected + item.min).toLocaleString("en-US", {style: "currency", currency: "USD", minimumFractionDigits: 2});
	};

	$scope.Init = function() {

		console.log($routeParams.id);

		$scope.data = {};

		$http({
			method: 'POST',
			url: "j_crowdfund_info",
			headers: {'Content-Type': 'application/x-www-form-urlencoded'},
			transformRequest: function(obj) {
				var str=[]; for(var p in obj) str.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]));
				return str.join("&");
			},
			data: {
				id: $routeParams.id
			}
		})
		.success(function(data) {

			console.log(data);

			if (data.res == 0) {

				$scope.data = data.data;

				$scope.item = {
					id: $routeParams.id,
					cover: data.data.cover,
					name: data.data.name,
					goal: data.data.goal,
					collected: data.data.сollected,
					min: 0
				};

				$scope.cover       = data.data.cover;
				$scope.embedLink   = data.data.video;
				$scope.next        = data.data.next;
				$scope.prev        = data.data.prev;
				$document[0].title = data.data.name;
			}
		});
	};

	$scope.Touch = function() {
		$location.path("/"+User	.id);
	}
});
