var nickname = "v1per14";
var channel = "";
var host = "http://localhost:8081"

// nickname = prompt('Введите никнем');	



$(document).ready(function () {
	

	$('ul.rooms-list li').on('click', function () {
		$('ul.rooms-list li').removeClass('active');
		
		$current = $(this);

		$current.addClass('active');

		$('#form_connect_to_room #channel_message').html('Вы выбрали канал <b>' + $current.text() + '</b>')
		$('#form_connect_to_room #room_field').val($current.data('value'))
		channel = $current.data('value');
	});

	$('#form_create_room').on('submit', function (e) {
		e.preventDefault

		$current = $(this);

		data = $current.serialize();

		$.post( host + "/room/create", data, function (data) {
			//Must return json
		})
	});

	$('#form_connect_to_room').on('submit', function (e) {
		e.preventDefault

		$current = $(this)

		data = $current.serialize();

		$.post( host + "/user/connect", data, function (data) {
			//Must return json
		})

		// var ws = new WebSocket('ws://localhost:8081');

		// ws.onopen = function () {
		// 	alert('connected');
		// }

		// ws.onerror = function () {
		// 	alert('error');
		// }	
	});
})