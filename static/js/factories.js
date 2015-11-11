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
})();