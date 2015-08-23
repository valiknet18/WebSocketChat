var user = {nickname: "quest" + Math.floor(Math.random() * 10), userHash: null, roomHash: null};
var host = "http://localhost:8081";

// nickname = prompt('Введите никнем');	

$(document).ready(function () {
	
	//Выбор канала (подсвечивает выбраный)
	$(document).on('click', 'ul.rooms-list li', function () {

		$('ul.rooms-list li').removeClass('active');
		
		$current = $(this);

		$current.addClass('active');

		$('#form_connect_to_room #channel_message').html('Вы выбрали канал <b>' + $current.text() + '</b>')
		$('#form_connect_to_room #room_field').val($current.data('value'))
		user.roomHash = $current.data('value');
	});

	//Вытягивает комнаты
	$.get(host + '/room/get', null, function (data) {
		result = "";

		for (room in data) {
			if (data[room].hash == user.roomHash) {
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
				if (data[room].hash == user.roomHash) {
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

	//Форма для создания канала
	$(document).on('submit', '#form_create_room', function (e) {
		e.preventDefault()

		nickname = $('#nickname_field').val();

		$current = $(this);

		data = $current.serialize();

		$.post( host + "/room/create", data, function (data) {
			result = "";
				
			for (room in data) {
				if (data[room].hash == user.roomHash) {
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

	//Форма которая подключается к комнате
	$('#form_connect_to_room').on('submit', function (e) {
		e.preventDefault()
		//ALGO
		//User write nickname, script send this nickname on server, generate user "object", and return hash (user_key), this hash save in variable, and send him with ws every time, when websocket will send.

		clearInterval(interval)

		$current = $(this);

		data = $current.serializeArray();

		console.log(data)

		if ($('ul li.active').length) {
			var ws;

			$.post(host + "/user/connect", {nickname: data[0]['value'], roomHash: data[1]['value']}, function (returnedData) {
				user.nickname = data[0]['value'];
				user.userHash = returnedData

				ws = new WebSocket('ws://localhost:8081/ws/' + user.userHash + '/connect');

				ws.onopen = function () {
					renderChat();
				}

				ws.onerror = function () {
					alert('error');
				}

				ws.onmessage = function (data) {
					console.log(data)
				}
			});
				
			
			//Event that send message
			$(document).on('submit', '#formSendMessage', function (e) {		
				e.preventDefault();

				data = $(this).serializeArray()

				console.log(data)

				ws.send({userHash: users.userHash, message: data[0]['value']});
			});	
		} else {
			alert('Нужно выбрать комнату')
		}
	});
})

//Рендерит чат
function renderChat() {
	users = ["user1", "user2"];

	someTestMessages = [{nickname: 'v1per14', message: 'Hello world'}, {nickname: 'Vasya', message: 'Hello my friend'}]

	$.get(host + "/room/users/" + user.roomHash, null, function (data) {
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
		$liMessage = $('<li>')
						.text(someTestMessages[key].nickname + ": " + someTestMessages[key].message)
						.appendTo($ulMessages);
	}

	$ulMessages.appendTo($leftDivMessageContent);

	$leftDivMessageContent.appendTo($leftDiv);

	$leftDivMessageFormContent	= $('<div>')
									.addClass('col-md-12')

	$leftDivMessageForm = $('<form id="formSendMessage">')
								

	$textareaDiv = $('<div>')
						.addClass('col-md-8')
						.addClass('form-group')
						.append('<textarea class="col-md-12 form-control" name="message">')
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