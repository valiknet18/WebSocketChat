var nickname = "v1per14";
var channel = "";
var host = "http://localhost:8081"

// nickname = prompt('Введите никнем');	

$(document).ready(function () {
	

	$(document).on('click', 'ul.rooms-list li', function () {

		$('ul.rooms-list li').removeClass('active');
		
		$current = $(this);

		$current.addClass('active');

		$('#form_connect_to_room #channel_message').html('Вы выбрали канал <b>' + $current.text() + '</b>')
		$('#form_connect_to_room #room_field').val($current.data('value'))
		channel = $current.data('value');
	});

	$.get(host + '/room/get', null, function (data) {
		result = "";

		for (room in data) {
			if (data[room].hash == channel) {
				result += '<li data-value="' + data[room].hash + '" class="active">' + data[room].name + '</li>'
			} else {
				result += '<li data-value="' + data[room].hash + '">' + data[room].name + '</li>'	
			}
		}	

		if (result == "") {
			result = "На данный момент не создано ни одного канала"
		}

		$('.rooms-list').html(result);
	});

	interval = setInterval(function() {
		$.get(host + '/room/get', null, function (data) {
			result = "";
				
			for (room in data) {
				if (data[room].hash == channel) {
					result += '<li data-value="' + data[room].hash + '" class="active">' + data[room].name + '</li>'
				} else {
					result += '<li data-value="' + data[room].hash + '">' + data[room].name + '</li>'	
				}
			}	

			if (result == "") {
				result = "На данный момент не создано ни одного канала"
			}

			$('.rooms-list').html(result);
		});
	}, 10000)

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