////
// app.js
////

var initMiUi = function () {
	console.log("PicoUI init");
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
					$('#header').prepend('<button class="button" id="' + msg.attributes.previd + '">' + msg.attributes.prevtxt + '</button>')
					$('#' + msg.attributes.previd).click(function (o) {
						$.get('/click?eid=' + $(this).attr('id'), {}, function (r) {});
					});
				}

				// Set the title
				$('#title').append(msg.attributes.title)
			} else if (msg.cmd === 'addbutton') {
				var classAttributes = "";
				if (msg.attributes.classAttr != null) {
					$.each(msg.attributes.classAttr, function (index, value) {
						classAttributes += " " + value; 
					});
				}
				
				var icon = "";
				if (msg.attributes.icon != "") {
					icon = '<i class="icon ' + msg.attributes.icon + '"></i>';
				}

				$('<button class="button ' + classAttributes + '" id="' + msg.attributes.eid + '">' + icon + msg.attributes.txt + '</button>').insertBefore(BEFORE);
				
				// Add the click handler
				$('#' + msg.attributes.eid).click(function (o) {
					$.get('/click?eid=' + $(this).attr('id'), {}, function (r) {});
				});
			} else if (msg.cmd === 'updateClassAttr') {
				var elem = $('#' + msg.attributes.eid);

				elem.removeAttr('class');
				elem.addClass('button');
				if (msg.attributes.classAttr != null) {
					$.each(msg.attributes.classAttr, function (index, value) {
						elem.addClass(value);
					});
				}
			} else if (msg.cmd === 'setIcon') {
				var icon = '<i class="icon ' + msg.attributes.icon + '"></i>';
				$('#' + msg.attributes.eid).children("i").remove();
				$('#' + msg.attributes.eid).prepend(icon);
			} else if (msg.cmd === 'addelement') {
				$('<' + msg.attributes.e + ' id="' + msg.attributes.eid + '">' + msg.attributes.txt + '</' + msg.attributes.e + ">").insertBefore(BEFORE);
			} else if (msg.cmd === 'updateinner') {
				$('#' + msg.attributes.eid).text(msg.attributes.txt);
			} else if (msg.cmd === 'addlist') {
				// $('<ul class="list" id="' + msg.attributes.eid + '"></ul>').insertBefore(BEFORE);
				$('<div class="list" id="' + msg.attributes.eid + '"></div>').insertBefore(BEFORE);
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
					itemHtml = '<a class="item" href="#" id="' + msg.attributes.eid + '">' + msg.attributes.txt + '</a>';
					console.log(itemHtml);
					console.log(msg.attributes.pid);
					$('#' + msg.attributes.pid).append(itemHtml);

					// Click event
					$('#' + msg.attributes.eid).click(function (o) {
						$.get('/click?eid=' + $(this).attr('id'), {}, function (r) {});
					});
				}
			} else if (msg.cmd === 'addtoggleitem') {
				//var toggleHtml = '<label class="toggle" id="' + msg.attributes.tid + '"> <input type="checkbox"><div class="track"><div class="handle"></div></div></label>';
				var toggleHtml = '<label class="toggle"> <input type="checkbox" id="' + msg.attributes.tid + '"><div class="track"><div class="handle"></div></div></label>';
				var itemHtml = '<div class="item item-toggle" id="' + msg.attributes.eid + '">' +
								msg.attributes.txt + toggleHtml + '</div>';

				$('#' + msg.attributes.pid).append(itemHtml);

				// Toggle event
				// document.querySelector('#' + msg.attributes.tid)
				// 	.addEventListener('toggle', function(event) {
				// 		console.log("Toggled!!");
    //                 	$.get('/toggle?eid=' + $(this).attr('id') +'&v=' + event["detail"]["isActive"]);
    //               	});

                $('#' + msg.attributes.tid).change(function () {
                	//console.log($(this).is(":checked"));
                	var isChecked = $(this).is(':checked') ? true : false;
                	//console.log(isChecked);

                	$.get('/toggle?eid=' + $(this).attr('id') + '&v=' + isChecked, {}, function (r) {});
                });
			} else if (msg.cmd === 'addinput') {
				$('<input id="' + msg.attributes.eid + '" type="' + msg.attributes.type + '" placeholder="' + msg.attributes.placeholder + '">').insertBefore(BEFORE);
			} else if (msg.cmd === 'getinput') {
				$.get('/state?msg=' + $('#' + msg.attributes.eid).val(), {}, function (r) {});
			}
		}
		setTimeout(function () {poll();}, 0);
	});
};