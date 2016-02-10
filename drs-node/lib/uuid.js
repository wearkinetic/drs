'use strict';

Object.defineProperty(exports, "__esModule", {
	value: true
});

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function build(seed) {
	var result = seed.toString(16);
	result = new Array(17 - result.length).join('0') + result;

	for (var i = 0; i < 16; i++) {
		result += Math.floor(Math.random() * 10).toString(16);
	}
	return result;
}

var MAX = Math.pow(2, 42);

var UUID = function () {
	function UUID() {
		_classCallCheck(this, UUID);
	}

	_createClass(UUID, [{
		key: 'ascending',
		value: function ascending() {
			return build(new Date().getTime());
		}
	}, {
		key: 'descending',
		value: function descending() {
			return build(MAX - new Date().getTime());
		}
	}]);

	return UUID;
}();

exports.default = new UUID();