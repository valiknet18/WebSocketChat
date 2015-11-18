(function () {
	var app = angular.module('chat', ['ngRoute', 'controllers', 'factories', 'ngWebSocket', 'LocalStorageModule']);

	app.config(['$routeProvider', function ($routeProvider) {
		$routeProvider
			.when("/", {
				controller: "HomeCtrl",
				templateUrl: "static/js/templates/index.html"
			})
			.when("/room/:slug", {
				controller: "RoomCtrl",
				templateUrl: "static/js/templates/room.html"
			})
			.otherwise({
				redirectTo: "/"
			});
	}]);
})();