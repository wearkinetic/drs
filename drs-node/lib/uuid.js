'use strict';

Object.defineProperty(exports, "__esModule", {
	value: true
});

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

var lastPushTime = 0;
var PUSH_CHARS = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz';
var PUSH_CHARS_REVERSE = PUSH_CHARS.split('').reverse().join('');
var lastRandChars = [];

var UUID = function () {
	function UUID() {
		_classCallCheck(this, UUID);
	}

	_createClass(UUID, [{
		key: 'now',
		value: function now() {
			return new Date().getTime();
		}
	}, {
		key: 'ascending',
		value: function ascending() {
			return generate(this.now(), PUSH_CHARS);
		}
	}, {
		key: 'descending',
		value: function descending() {
			return generate(this.now(), PUSH_CHARS_REVERSE);
		}
	}]);

	return UUID;
}();

function generate(time, pool) {
	var len = pool.length;
	var now = time;
	var duplicateTime = now === lastPushTime;
	lastPushTime = now;

	var timeStampChars = new Array(8);
	for (var i = 7; i >= 0; i--) {
		timeStampChars[i] = pool.charAt(now % len);
		// NOTE: Can't use << here because javascript will convert to int and lose the upper bits.
		now = Math.floor(now / len);
	}
	if (now !== 0) throw new Error('We should have converted the entire timestamp.');

	var id = timeStampChars.join('');

	if (!duplicateTime) {
		for (var _i = 0; _i < 12; _i++) {
			lastRandChars[_i] = Math.floor(Math.random() * len);
		}
	} else {
		var _i2 = void 0;
		for (_i2 = 11; _i2 >= 0 && lastRandChars[_i2] === len - 1; _i2--) {
			lastRandChars[_i2] = 0;
		}
		lastRandChars[_i2]++;
	}
	for (var _i3 = 0; _i3 < 12; _i3++) {
		id += pool.charAt(lastRandChars[_i3]);
	}
	if (id.length !== 20) throw new Error('Length should be 20.');
	return id;
}

exports.default = new UUID();