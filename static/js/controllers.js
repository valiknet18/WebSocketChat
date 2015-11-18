(function () {
	app = angular.module('controllers', ["ngRoute"]);

	app.controller('HomeCtrl', ['$scope', '$location', 'createUser', 'localStorageService', 'getRooms', function ($scope, $location, createUser, localStorageService, getRooms) {
		$scope.user = {}
		$scope.room = {}

		$scope.createNickname = function () {
			if (!$scope.user.nickname) {
				alert("Спочатку треба ввести нікнейм");
			} else {
				createUser($scope.user.nickname)
				.then(function (response) {
					localStorageService.set("user_token", response.data.Hash)

					$scope.channels = getRooms().then(function (response) {
						$scope.channels = response.data;

						$element = angular.element("#channels-block")

						$element.fadeIn("slow");
						$element.removeClass("hidden");

						angular.element(".list-group.list-group-item").on("click", function () {
							console.log($(this).data("slug"));
						})
					}, function (error) {

					});

				}, function (error) {
					console.log("Користувач не створений")
				});
			}
		}

		$scope.connectToRoom = function () {
			if (!$scope.room.slug) {
				alert("Для використання потрібно обов’язково ввести свій нікнейм")
			} else {
				console.log($scope.room.slug)

				$location.path("/room/" + $scope.room.slug)
			}
		}

		$scope.createNewRoom = function () {
			if (!$scope.room.name) {
				alert("Поле назва кімнати обов’язково до заповнення!")
			} else {
				console.log("Created")
				$element = angular.element("#myModal")

				$element.modal('hide');
			}
		}	
	}]);

	app.controller("RoomCtrl", ["$scope", "$routeParams", "connectToRoom", "localStorageService", function ($scope, $routeParams, connectToRoom, localStorageService) {
		$scope.message = {
			maxChars: 250,
			currentChars: 250,
		}

		$scope.connect = connectToRoom($routeParams.slug, localStorageService.get("user_token"));
	}]);
})();