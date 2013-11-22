////
// app.js
////

var initMiUi = function () {
	console.log("MiUi init");
	$.get("/init", {}, function (r) {
		poll();
	});
};

var dispatch = function (msg) {
	msg_json = JSON.parse(msg);
	if (msg_json.cmd === "timeout") {
		return null;
	} else if (msg_json.cmd === "newpage") {
		console.log("newpage");

		// Clear the page
		$('#header').html('<h1 class="title" id="title"></h1>');
		$('#padded').html('<p id="end"></p>');

		return null;
	} else {
		return msg_json;
	}
};

var BEFORE = "#end";
function poll() {
	$.get("/poll", {}, function (m) {
		msg = dispatch(m);
		if (msg != null) {
			// A new page was posted
			if (msg.cmd === 'pagepost') {
				// Set the back button, if specified
				if (msg.attributes.previd) {
					$('#header').prepend('<a href="#" id="' + msg.attributes.previd + '" class="button-prev">' + msg.attributes.prevtxt + '</a>');
					$('#' + msg.attributes.previd).click(function (o) {
						$.get('/click?eid=' + $(this).attr('id'), {}, function (r) {});
					});
				}

				// Set the title
				$('#title').append(msg.attributes.title)
			} else if (msg.cmd === 'addbutton') {
				$('<a class="button" id="' + msg.attributes.eid + '">' + msg.attributes.txt + '</a>').insertBefore(BEFORE);
				$('#' + msg.attributes.eid).click(function (o) {
					$.get('/click?eid=' + $(this).attr('id'), {}, function (r) {});
				});
			} else if (msg.cmd === 'addelement') {
				$('<' + msg.attributes.e + ' id="' + msg.attributes.eid + '">' + msg.attributes.txt + '</' + msg.attributes.e + ">").insertBefore(BEFORE);
			} else if (msg.cmd === 'updateinner') {
				$('#' + msg.attributes.eid).text(msg.attributes.txt);
			} else if (msg.cmd === 'addlist') {
				$('<ul class="list" id="' + msg.attributes.eid + '"></ul>').insertBefore(BEFORE);
			} else if (msg.cmd === 'addlistitem') {
				var chevron = "";
				if (msg.attributes.chevron) {
					chevron = '<span class="chevron"></span>';
				}

				var itemHtml = "";
				if (msg.attributes.toggle) {
					var toggle = '<div class="toggle" id="' + msg.attributes.tid +
								'"><div class="toggle-handle"></div></div>';

					itemHtml = '<li id="' + msg.attributes.eid + '"><a>' + msg.attributes.txt + toggle + '</a></li>';
					$('#' + msg.attributes.pid).append(itemHtml);

					// Toggle event
					if (msg.attributes.toggle) {
						document.querySelector('#' + msg.attributes.tid)
							.addEventListener('toggle', function(event) {
								console.log("Toggled!!");
                    			$.get('/toggle?eid=' + $(this).attr('id') +'&v=' + event["detail"]["isActive"]);
                  			});
                  	}
				} else {
					itemHtml = '<li id="' + msg.attributes.eid + '"><a>' + msg.attributes.txt + chevron + '</a></li>';
					$('#' + msg.attributes.pid).append(itemHtml);	

					// Click event
					$('#' + msg.attributes.eid).click(function (o) {
						$.get('/click?eid=' + $(this).attr('id'), {}, function (r) {});
					});
				}
			} else if (msg.cmd === 'addinput') {
				$('<input id="' + msg.attributes.eid + '" type="' + msg.attributes.type + '" placeholder="' + msg.attributes.placeholder + '">').insertBefore(BEFORE);
			} else if (msg.cmd === 'getinput') {
				$.get('/state?msg=' + $('#' + msg.attributes.eid).val(), {}, function (r) {});
			}
		}
		setTimeout(function () {poll();}, 0);
	});
};