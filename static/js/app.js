var nickname = "v1per14";
var channel = "";
var host = "http://localhost:8081";

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
	}, 10000);

	$(document).on('submit', '#form_create_room', function (e) {
		e.preventDefault()

		nickname = $('#nickname_field').val();

		$current = $(this);

		data = $current.serialize();

		$.post( host + "/room/create", data, function (data) {
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

			$('#form_create_room #roomname_field').val('')
		})
	});	

	$('#form_connect_to_room').on('submit', function (e) {
		e.preventDefault()

		clearInterval(interval)

		$current = $(this)

		data = $current.serialize();

		var ws = new WebSocket('ws://localhost:8081/room/' + room);

		ws.onopen = function () {
			ws.send(JSON.stringify({nickname: nickname}))
			
			renderChat();
		}

		ws.onerror = function () {
			alert('error');
		}	
		
		//Event that send message
		$(document).on('submit', '#form-send-message', function (e) {		
			e.preventDefault();

			ws.send()
		});
	});
})

function renderChat() {
	users = ["user1", "user2"];

	someTestMessages = [{nickname: v1per14, message: 'Hello world'}, {nickname: 'Vasya', message: 'Hello my friend'}]

	$.get(host + "/room/users/" + room, null, function (data) {
		users = data;
	});

	$divMain = $('<div>')
					.addClass('col-md-12');
	$pageHeaderDiv = $('<div>')
						.addClass('page-header')
						.appendTo($divMain);

	$leftDiv = $('<div>')
					.addClass('col-md-8');

	$leftDivMessageContent = $('<div>')
								.addClass('col-md-12')

	$ulMessages = $('<ul>');

	for (key in someTestMessages) {
		$liMessage = $('<li>');
	}

	$leftDivMessageContent.appendTo($leftDiv);

	$leftDivMessageFormContent	= $('<div>')
									.addClass('col-md-12')

	$leftDivMessageForm = $('<form>')
								

	$textareaDiv = $('<div>')
						.addClass('col-md-8')
						.addClass('form-group')
						.append('<textarea class="col-md-12 form-control">')
						.appendTo($leftDivMessageForm);

	$buttonDiv = $('<div>')
						.addClass('col-md-4')
						.addClass('form-group');

	$buttonSendMessage = $('<button type="submit">')
								.text('Send message')
								.addClass('btn btn-success')
								.appendTo($buttonDiv);

	$buttonDiv.appendTo($leftDivMessageForm)							
	$leftDivMessageForm.appendTo($leftDivMessageContent);

	$leftDivMessageFormContent.appendTo($leftDiv);

	$leftDiv.appendTo($divMain);
	$rightDiv = $('<div>')
		.addClass('col-md-4')
		.appendTo($divMain);

	$ul = $('<ul id="users">');

	for (key in users) {
		$li = $('<li>').text(users[key]);
		$li.appendTo($ul);
	}

	$rightDiv.append($ul);

	// $divMain
	// 	.append($leftDiv)
	// 	.append($rightDiv);

	console.log($divMain)

	$('body > .container').html($divMain);


	// html = '<div class="col-md-12"><div class="page-header"></div>';

	// html += "<div class='col-md-8'><div class='col-md-12'></div><div class='col-md-12'><form id='form-send-message'><div><textarea></textarea></div><div><button type='submit' class='btn btn-success'>Отправить</button></div></form></div></div>"

	// html += "<div id='users' class='col-md-4'><ul>"

	// ul = "";

	// for (key in users) {
	// 	ul += '<li>' + users[key] + '</li>';
	// }

	// html += ul + "</ul>";

	// html += "</div></div>";

	// $('body > .container').html(html);
}