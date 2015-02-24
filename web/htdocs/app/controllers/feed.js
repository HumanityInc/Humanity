
app.controller('FeedController', function($scope, $http, $timeout, $document, $location) {

	// console.log(myData);
	
	$scope.showFeed = true;
	$scope.showCreate = false;
	$scope.showMenu = true;
	$scope.showCrowdfund = false;

	$scope.feed = [];
	$scope.feedX2Count = 0;
	$scope.feedX3Count = 0;
	$scope.step = 25;

	$scope.search = {placeholder: "find (person, topic, interest, song, movie, place, biz, etc.)"};

	$scope.form = {
		name: "",
		goal: "",
		video: "", 
		cover: ""
	};

	$scope.avatar = {
		"background": "url("+User.picture+") 50% 50% no-repeat",
		"background-size": "cover"
	};

	$scope.avatarUpload = {
		// url: 'http://127.0.0.1:1991/upload',
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
						var str=[];
						for(var p in obj) str.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]));
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

	$scope.fileUpload = {
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
				$scope.form.cover = res.cover;
				$scope.loading = false;
			},
			error: function(uploader, error) {

				$scope.loading = false;
				alert(error.message);
			}
		}
	};

	// ==

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

	// ==

	$scope.SearchFocus = function() {
		$scope.search.placeholder = "";
	};
	$scope.SearchBlure = function() {
		$scope.search.placeholder = "find (person, topic, interest, song, movie, place, biz, etc.)";
	};

	// ==

	$scope.ShowCreate = function() {
		$scope.showCreate = true;
	};
	$scope.HideCreate = function() {
		$scope.showCreate = false;
	};


	$scope.Create = function() {

		// alert($scope.form.cover)

		$http({
			method: 'POST',
			url: "j_crowdfund",
			headers: {'Content-Type': 'application/x-www-form-urlencoded'},
			transformRequest: function(obj) {
				var str=[];
				for(var p in obj) str.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]));
				return str.join("&");
			},
			data: {
				goal:  $scope.form.goal,
				name:  $scope.form.name,
				video: $scope.form.video,
				cover: $scope.form.cover
			}
		})
		.success(function(data) {

			if (data.res == 0) {
				// $scope.form.crowdfundId = data.data.id;
				$scope.HideCreate();
				$scope.Load();
			}
		});
	};

	$scope.Load = function() {

		$http({
			method: 'POST',
			url: "j_feed",
			headers: {'Content-Type': 'application/x-www-form-urlencoded'},
			transformRequest: function(obj) {
				var str=[];
				for(var p in obj) str.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]));
				return str.join("&");
			},
			data: {
			}
		})
		.success(function(data) {
			if (data.res == 0) {
				
				var feed = [], X2Count = 0, X3Count = 0;

				var i=0, l=data.data.length, type=2, buf=[];

				while(i < l) {

					var item = data.data[i];

					var cf = {
						id: item.id,
						cover: item.cover,
						goal: item.goal,
						collected: item.сollected,
						min: 0
					};

					buf.push(cf);

					if (buf.length == type) {

						feed.push({type: type, items: buf});
						buf = [];

						// if (type == 2) type = 3;
						// else type = 2;
					}

					i++;
				}

				if (buf.length > 0) {
					feed.push({type: type, items: buf});
				}

				for(var i=0, l=feed.length; i<l; i++) {
					switch(feed[i].type) {
					case 2: X2Count++; break;
					case 3: X3Count++; break;
					}
				}

				console.log("Load", feed);

				$scope.feedX2Count = X2Count;
				$scope.feedX3Count = X3Count;
				$scope.feed = feed;

				if ($scope.OnResize) {
					$scope.OnResize($scope.getWindowDimensions());
				}
			}
		});
	};

	$scope.Init = function() {

		$scope.Load();
		$document[0].title = "Humanity - Made In Ukraine!";
	};

	var raw = $document[0].body;
		
	$scope.Wheel = function(event, delta, deltaX, deltaY) {

		raw.scrollLeft -= delta;
		event.preventDefault();
	}

	$scope.GotoCrowdfund = function(item) {
		$location.path("/crowdfund/"+item.id);
	};
});
