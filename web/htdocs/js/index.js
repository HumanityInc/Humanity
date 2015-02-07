
var humanity = humanity || {};
humanity.ui = {};
humanity.model = {};

// var elem = document.getElementById("myvideo");
// if (elem.requestFullscreen) {
// elem.requestFullscreen();
// } else if (elem.mozRequestFullScreen) {
// elem.mozRequestFullScreen();
// } else if (elem.webkitRequestFullscreen) {
// elem.webkitRequestFullscreen();
// }

function createCookie(name, value, days) {
	if(days) {
		var date = new Date();
		date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
		var expires = "; expires=" + date.toGMTString();
	}
	else var expires = "";
	document.cookie = name + "=" + value + expires + "; path=/";
}

function getCookie(c_name) {
	if(document.cookie.length > 0) {
		c_start = document.cookie.indexOf(c_name + "=");
		if(c_start != -1) {
			c_start = c_start + c_name.length + 1;
			c_end = document.cookie.indexOf(";", c_start);
			if(c_end == -1) {
				c_end = document.cookie.length;
			}
			return unescape(document.cookie.substring(c_start, c_end));
		}
	}
	return "";
}

function detectmob() { 
	if( navigator.userAgent.match(/Android/i)
		|| navigator.userAgent.match(/webOS/i)
		|| navigator.userAgent.match(/iPhone/i)
		|| navigator.userAgent.match(/iPad/i)
		|| navigator.userAgent.match(/iPod/i)
		|| navigator.userAgent.match(/BlackBerry/i)
		|| navigator.userAgent.match(/Windows Phone/i)
	){
		return true;
	} else {
		return false;
	}
}

$(function() {
	'use strict';

	if (detectmob()) {
		$('video').remove();
	}

	var $signup_create = $('#signup_create'),
		$signup_main = $('#signup_main'),
		$login_box = $('.login_box'),
		$success2 = $('#success2'),
		$success = $('#success'),
		$signin = $('#signin'),
		$email = $('#email'),
		$main = $('#main'),
		$all = $signin.add($signup_main).add($signup_create).add($main).add($success).add($email).add($success2);

	var $win = $(window);
	$win.resize(function() {
		if ($win.height() > 900) {
			$login_box.removeClass('small');
		} else {
			$login_box.addClass('small');
		}
	}).resize();

	$login_box.show();

	$('input').focus(function(){
		$(this).data('placeholder',$(this).attr('placeholder'))
		$(this).attr('placeholder','');
	})
	.blur(function(){
		$(this).attr('placeholder',$(this).data('placeholder'));
	});

	var $video = $("video");

	if (getCookie("muted")|0) {
		$video.prop('muted', true);
		$('.volume').addClass('off');
	}

	$('.volume').click(function() {
		if ($(this).hasClass('off')) {
			$video.prop('muted', false);
			$(this).removeClass('off');
			createCookie("muted", "0", 7);
		} else {
			$video.prop('muted', true);
			$(this).addClass('off');
			createCookie("muted", "1", 7);
		}
		
	});

	var $color = $('.color');
	$color.mouseenter(function() {

		var c = $(this).data('color');
		$login_box.addClass('s'+c);
	})
	.mouseleave(function() {

		var c = $(this).data('color');
		$login_box.removeClass('s'+c);
	})
	.click(function() {

		var c = $(this).data('color'),
			cur = $login_box.data('cur');

		$login_box.removeClass('f'+cur).data('cur', c).addClass('f'+c);

		$color.removeClass('active');
		$(this).addClass('active');
	});

	$("#save").click(function() {

		var $parent = $(this).parent().parent();
		$parent.find('.wrong').hide();

		var email = $('input[name="email"]', $parent);

		if (email.val() == "") {
			email.focus().next().show();
			return;
		}

		$.ajax({
			url: "/j_setemail",
			type: "POST",
			dataType: "JSON",
			data: {
				email: email.val()
			}
		})
		.done(function(data){

			if (data.res == 0) {
				document.location.href = "#!success"
			} else {
				switch(data.error) {
				case "INVALID_EMAIL":
					email.focus().next().show();
					break;
				}
			}
		});
	});

	$("#login").click(function() {

		var $parent = $(this).parent().parent();

		$parent.find('.wrong').hide()

		var email = $('input[name="email"]', $parent),
			password = $('input[name="password"]', $parent)

		if (email.val() == "") {
			email.focus().next().show();
			return;
		}
		if (password.val() == "") {
			password.focus().next().show();
			return;
		}

		$.ajax({
			url: "/j_login",
			type: "POST",
			dataType: "JSON",
			data: {
				email: email.val(),
				password: password.val()
			}
		})
		.done(function(data){

			if (data.res == 0) {
				document.location.href = "#!login"
			} else {
				switch(data.error) {
				case "INVALID_EMAIL":
					email.focus().next().show();
					break;
				case "INVALID_EMAIL_OR_PASSWORD":
					password.focus().next().show();
					break;
				}
			}
		});
	});

	$("#create1, #create2").click(function() {

		var $parent = $(this).parent().parent();

		$parent.find('.wrong').hide()

		var email = $('input[name="email"]', $parent),
			password = $('input[name="password"]', $parent),
			password2 = $('input[name="password2"]', $parent),
			first_name = $('input[name="first_name"]', $parent),
			last_name = $('input[name="last_name"]', $parent);
		
		if (email.val() == "") {
			email.focus().next().show();
			return;
		}
		if (password.val() == "") {
			password.focus().next().show();
			return;
		}
		if (password2.val() == "") {
			password2.focus().next().show();
			return;
		}
		if (first_name.val() == "") {
			first_name.focus().next().show();
			return;
		}
		if (last_name.val() == "") {
			last_name.focus().next().show();
			return;
		}

		$.ajax({
			url: "/j_register",
			type: "POST",
			dataType: "JSON",
			data: {
				email: email.val(),
				password: password.val(),
				password2: password2.val(),
				last_name: last_name.val(),
				first_name: first_name.val()
			}
		})
		.done(function(data){

			if (data.res == 0) {
				document.location.href = "#!success"
			} else {
				switch(data.error) {
				case "EMAIL_ALREADY_EXISTS":
				case "INVALID_EMAIL":
					email.focus().next().show();
					break;
				case "PASSWORDS_NOT_EQUAL":
					password2.focus().next().show();
					break;
				}
			}
		});

	});

	var Router = Backbone.Router.extend({
		routes: {
			"": "main",
			"!": "main",
			"!main": "main",
			"!email": "email",
			"!signin": "signin",
			"!signin/create": "create_signup",
			"!main/signup": "main_signup",
			"!success": "success",
			"!login": "login"
		},
		initialize: function() {
			var stat = Backbone.history.start({pushState: false, root: '/'});
		}, 
		main: function() {
			$all.not($main).hide();
			$main.show();
		},
		signin: function() {
			$all.not($signin).hide();
			$signin.show().find('input');
		},
		main_signup: function() {
			$all.not($signup_main).hide();
			$signup_main.show().find('input');	
		},
		create_signup: function() {

			var email = $('input[name="email"]', $signin),
				password = $('input[name="password"]', $signin);
			
			$('input[name="email"]', $signup_create).val(email.val());
			$('input[name="password"]', $signup_create).val(password.val());

			$all.not($signup_create).hide();
			$signup_create.show().find('input');
		},
		login: function() {
			$all.not($success2).hide();
			$success2.show();
		},
		success: function() {
			$all.not($success).hide();
			$success.show();
		},
		email: function() {
			$all.not($email).hide();
			$email.show().find('input');
		}
	});

	humanity.router = new Router();
});