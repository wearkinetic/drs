'use strict';

Object.defineProperty(exports, "__esModule", {
	value: true
});

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

var _querystring = require('querystring');

var _querystring2 = _interopRequireDefault(_querystring);

var _pipe = require('../pipe');

var _pipe2 = _interopRequireDefault(_pipe);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function _possibleConstructorReturn(self, call) { if (!self) { throw new ReferenceError("this hasn't been initialised - super() hasn't been called"); } return call && (typeof call === "object" || typeof call === "function") ? call : self; }

function _inherits(subClass, superClass) { if (typeof superClass !== "function" && superClass !== null) { throw new TypeError("Super expression must either be null or a function, not " + typeof superClass); } subClass.prototype = Object.create(superClass && superClass.prototype, { constructor: { value: subClass, enumerable: false, writable: true, configurable: true } }); if (superClass) Object.setPrototypeOf ? Object.setPrototypeOf(subClass, superClass) : subClass.__proto__ = superClass; }

var WS = function (_Pipe) {
	_inherits(WS, _Pipe);

	function WS() {
		var query = arguments.length <= 0 || arguments[0] === undefined ? {} : arguments[0];
		var proto = arguments.length <= 1 || arguments[1] === undefined ? 'ws' : arguments[1];

		_classCallCheck(this, WS);

		var _this = _possibleConstructorReturn(this, Object.getPrototypeOf(WS).call(this));

		_this._query = query;
		_this._proto = proto;
		return _this;
	}

	_createClass(WS, [{
		key: '_connect',
		value: function _connect(host) {
			var _this2 = this;

			return new Promise(function (resolve, reject) {
				var ws = new global.WebSocket(_this2._proto + '://' + host + '/socket?' + _querystring2.default.stringify(_this2._query));
				ws.onerror = reject;
				ws.onopen = function () {
					return resolve(session(ws));
				};
			});
		}
	}]);

	return WS;
}(_pipe2.default);

exports.default = WS;


function session(ws) {
	return {
		send: function send(data) {
			return new Promise(function (resolve) {
				ws.send(data);
				resolve();
			});
		},
		on: function on(action, cb) {
			if (action === 'data') ws.onmessage = function (msg) {
				return cb(msg.data);
			};
			if (action === 'close') ws.onclose = function () {
				return cb();
			};
		},
		close: function close() {
			ws.close();
		}
	};
}