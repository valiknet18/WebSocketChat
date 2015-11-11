(function () {
	app = angular.module('controllers', ["ngRoute"]);

	app.controller('HomeCtrl', ['$scope', '$location', function ($scope, $location) {
		$scope.user = {}
		$scope.room = {}

		$scope.createNickname = function () {
			if (!$scope.user.nickname) {
				alert("Спочатку треба ввести нікнейм");
			} else {
				$element = angular.element("#channels-block")

				$element.fadeIn("slow");
				$element.removeClass("hidden");

				angular.element(".list-group.list-group-item").on("click", function () {
					console.log($(this).data("slug"));
				})
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

	app.controller("RoomCtrl", ["$scope", "$routeParams", function ($scope, $routeParams) {
		$scope.message = {
			maxChars: 250,
			currentChars: 250,
		}
	}]);
})();