var user = {nickname: "quest" + Math.floor(Math.random() * 10), userHash: null, roomHash: null};

// local host
// var domain = "http://"
// var host = "localhost:8080";
// var wsDomain = "ws://"

// prod host
var domain = "https://"
var host = "go-test-app-project.herokuapp.com";
var wsDomain = "wss://"

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
	$.get(domain + host + '/room/get', null, function (data) {
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
		$.get(domain + host + '/room/get', null, function (data) {
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

		$.post(domain + host + "/room/create", data, function (data) {
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

	flagActive = false;

	if (!flagActive) {
		flagActive = true

		//Форма которая подключается к комнате
		$('#form_connect_to_room').on('submit', function (e) {
			e.preventDefault()
			//ALGO
			//User write nickname, script send this nickname on server, generate user "object", and return hash (user_key), this hash save in variable, and send him with ws every time, when websocket will send.

			clearInterval(interval)

			$current = $(this);

			data = $current.serializeArray();

			if ($('ul li.active').length) {
				$current.remove();

				var ws;

				$.post(domain + host + "/user/connect", {nickname: data[0]['value'], roomHash: data[1]['value']}, function (returnedData) {
					user.nickname = data[0]['value'];
					user.userHash = returnedData

					ws = new WebSocket(wsDomain + host + '/ws/' + user.userHash + '/connect');

					ws.onopen = function () {
						console.log("Соединение установлено")

						renderChat();
						ws.send('Hello world')
					}

					ws.onerror = function () {
						console.log("Соединение разорвано")
					}

					ws.onmessage = function (event) {
						console.log(event.data)

						renderMessage(event.data)
					}
				});
					
				
				//Event that send message
				$(document).on('submit', '#formSendMessage', function (e) {		
					e.preventDefault();

					data = $(this).serializeArray()

					if (data[0]['value']) {
						obj = {UserHash: user.userHash, Message: data[0]['value']}

						json = JSON.stringify(obj)	

						ws.send(json);	
						
						$('#formSendMessage textarea').val("")
					} else {
						alert("Нужно ввести сообщение")
					}
				});	
			} else {
				alert('Нужно выбрать комнату')
			}
		});
	}
	
})

function renderMessage (data) {
	userInfo = JSON.parse(data)

	var $mainLi = $('<li>')

	$message = $('<div class="col-md-12">')

	$messageText = $('<div class="col-md-8">')
		.html("<b>" + userInfo.user.nickname  + "</b>: " + userInfo.message)
		.appendTo($message)

	$messageDate = $('<div class="col-md-4">')
		.html("<span class='pull-right'>" + userInfo.created_at + "</span>")
		.appendTo($message)

	$message.appendTo($mainLi);

	$messageContent = $('#messageContent ul')
					
	$messageContent.append($mainLi)

	$messageContent.scrollTop($messageContent[0].scrollHeight - $messageContent.height());

	twemoji.size = '16x16';

	 // Parse the document body and
	 // insert <img> tags in place of Unicode Emojis
	 twemoji.parse(document.getElementById('messageContent'));	
}

//Рендерит чат
function renderChat() {
	users = [];

	$.get(domain + host + "/room/users/" + user.roomHash, null, function (data) {
		users = data;

		$('#users').html('')

		for (key in users.users) {
			$li = $('<li>').text(users.users[key].nickname);
			$li.appendTo($ul);
		}
	});

	setInterval(function () {
		$.get(domain + host + "/room/users/" + user.roomHash, null, function (data) {
			users = data;

			$('#users').html('')

			for (key in users.users) {
				$li = $('<li>').text(users.users[key].nickname);
				$li.appendTo($ul);
			}
		});
	}, 3000)

	$divMain = $('<div>')
					.addClass('col-md-12');
	$pageHeaderDiv = $('<div>')
						.addClass('page-header')
						.appendTo($divMain);

	$leftDiv = $('<div>')
					.addClass('col-md-8');

	$leftDivMessageContent = $('<div id="messageContent">')
								.addClass('col-md-12')

	$ulMessages = $('<ul class="col-md-12">');

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

	$('<h4>')
			.text('Список активных пользователей:')
			.appendTo($rightDiv)

	$rightDiv.appendTo($divMain);	

	$ul = $('<ul id="users">');

	$rightDiv.append($ul);

	console.log($divMain)

	$('body > .container').html($divMain);

	$('#formSendMessage textarea').emojiPicker({
		width: '150px',
		height: '300px'
	});
}