(function () {
	app = angular.module('factories', []);

	app.directive("clickRoom", function () {
		return function ($scope, $element) {
			$element.bind("click", function (data) {
				$("ul li").removeClass("active");
				$(this).addClass("active")

				$scope.room.slug = data.target.dataset.slug;
			});
		}
	});

	app.directive("onWriteMessage", function () {
		return function ($scope, $element) {
			$element.bind("keyup", function (data) {
				$scope.message.currentChars = $scope.message.maxChars - data.target.value.length
				// console.log($scope.message.currentChars);

				angular.element("#count-chars").text($scope.message.currentChars)
			});
		}
	});

	app.factory("createUser", ['$http', '$q', function ($http, $q) {
		var defer = $q.defer();

		return function (nickname) {
			$http
				.post("/users/create", {nickname: nickname})
				.then(function (response) {
					defer.resolve(response)
				}, function (error) {
					defer.reject(error)
				})
			;

			return defer.promise
		}
	}]);

	app.factory("getRooms", ['$http', '$q', function ($http, $q) {
		var defer = $q.defer()

		return function () {
			$http
				.get("/rooms/get")
				.then(function (response) {
					defer.resolve(response)
				}, function (error) {
					defer.reject(error)
				})
			;

			return defer.promise
		}
	}]);

	app.factory("createRoom", ['$http', '$q', 'avia_host', function ($http, $q, avia_host) {
		return function (roomName) {
			var defer = $q.defer();

			$http
				.post(avia_host + "/rooms/create", {room: roomName})
				.then(function (data) {
					defer.resolve(data)
				}, function (error) {
					defer.reject(error)
				})
			;

			return defer.promise
		}
	}]);

	app.factory("connectToRoom", ['$websocket', function ($websocket) {

		return function (roomHash, userHash) {
			var dataStream = $websocket("ws://127.0.0.1:8000/ws/" + roomHash + "/connect/" + userHash);

			var collection = [];

			dataStream.onMessage(function (message) {
				collection.push(message)
			});

			return {
				collection: collection,
				send: function (message) {
					data = {
						userHash: userHash,
						message: message 
					};

					$dataStream.send(data)
				}
			}
		}
	}]);

	app.factory("createRoom", ['$http', '$q', function ($http, $q) {
		var defer = $q.defer();

		return function (roomName) {
			$http
				.post("/rooms/create", {room: roomName})
				.then(function (response) {
					defer.resolve(response)
				}, function (error) {
					defer.reject(error)
				})
			;

			return defer.promise
		}	
	}]);
})();