var nickname = "v1per14";
// nickname = prompt('Введите никнем');	

$(document).ready(function () {
	var ws = new WebSocket('ws://localhost:8081');

	ws.onopen = function () {
		alert('connected');
	}

	ws.onerror = function () {
		alert('error');
	}
})